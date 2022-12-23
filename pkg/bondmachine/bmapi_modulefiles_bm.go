package bondmachine

const (
	moduleFilesBm = `#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/module.h>
#include <linux/kdev_t.h>
#include <linux/fs.h>
#include <linux/cdev.h>
#include <linux/device.h>
#include <linux/slab.h>    //kmalloc()
#include <linux/uaccess.h> //copy_to/from_user()
#include <linux/sysfs.h>
#include <linux/kobject.h>
#include <linux/interrupt.h>
#include <linux/irq.h>
#include <asm/io.h>
#include <linux/workqueue.h>

#define BMBASEADDR 0x43c00000
#define BMLEN 24

#define REGSIZE 1

#define BMI0_O 0
#define BMPS2PL_O 4
#define BMSTATES_O 8
#define BMO0_O 12
#define BMPL2PS_O 16
#define BMCHANGES_O 20

#define I0_ID (uint8_t)0

#define O0_ID (uint8_t)0

#define CH_I0_RECV 0x80000000
#define CH_O0 0x40000000
#define CH_O0_VALID 0x20000000

#define PL2PS_I0_RECV 0x80000000
#define PL2PS_O0_VALID 0x40000000

#define PS2PL_I0_VALID 0x80000000
#define PS2PL_O0_RECV 0x40000000

#define bm_reg_size 1
#define write_buff_size 1024

#define cmdNEWVAL (uint8_t)0   // 000 00000
#define cmdDVALIDH (uint8_t)32 // 001 00000
#define cmdDVALIDL (uint8_t)64 // 010 00000
#define cmdDRECVH (uint8_t)96  // 011 00000
#define cmdDRECVL (uint8_t)128 // 100 00000
#define cmdHANDSH (uint8_t)160 // 101 00000
#define cmdKEEP (uint8_t)192   // 110 00000
#define cmdMASK (uint8_t)224   // 111 00000

#define stateWAIT (uint8_t)0
#define stateHSRECV (uint8_t)1
#define stateMASKSENT (uint8_t)2
#define stateACK (uint8_t)3
#define stateCONNECT (uint8_t)4

void *bmptr;

dev_t dev = 0;
static struct class *dev_class;
static struct cdev bm_cdev;
uint8_t *write_buffer;
int i;

wait_queue_head_t wait_queue_bm;
int wait_queue_flag = 0; // 0 - Blocked, 1 - Reply from write op, 2 - Reply to keep, 3 - Actions from the interrupt, 4 - file open , 5 - Exit

uint8_t bmacc_state = stateWAIT;

uint8_t hs = 0;
uint8_t hsmask = (uint8_t)85;

uint8_t *replies;

uint32_t changes = 0;
uint32_t value = 0;

int seqlen = 0;

int waitvalue = 0;
uint8_t *value_buffer;
int regvalue = 0;

struct workqueue_struct *wq;

struct work_data
{
        struct work_struct work;
        int cs;
};

static void
work_handler(struct work_struct *work)
{
        struct work_data *data = (struct work_data *)work;
        wait_queue_flag = data->cs;
        kfree(data);
        //pr_info("%d\n", wait_queue_flag);
        wake_up_interruptible(&wait_queue_bm);
}

// Interrupt Request number
#define IRQ_NO 46

// Interrupt handler for IRQ 11.
static irqreturn_t irq_handler(int irq, void *dev_id)
{
        struct work_data *irqwork;
        if (bmacc_state == stateCONNECT)
        {
                irqwork = kmalloc(sizeof(struct work_data), GFP_KERNEL);
                INIT_WORK(&irqwork->work, work_handler);
                irqwork->cs = 3;
                queue_work(wq, &irqwork->work);
        }
        // pr_info("Shared IRQ: Interrupt Occurred");
        return IRQ_HANDLED;
}

/*
** Function Prototypes
*/
static int __init bmmod_init(void);
static void __exit bmmod_exit(void);
static int bm_open(struct inode *inode, struct file *file);
static int bm_release(struct inode *inode, struct file *file);
static ssize_t bm_read(struct file *filp, char __user *buf, size_t len, loff_t *off);
static ssize_t bm_write(struct file *filp, const char *buf, size_t len, loff_t *off);

/*
** File Operations structure
*/
static struct file_operations fops =
    {
        .owner = THIS_MODULE,
        .read = bm_read,
        .write = bm_write,
        .open = bm_open,
        .release = bm_release,
};

/*
** This function will be called when we open the Device file
*/
static int bm_open(struct inode *inode, struct file *file)
{
        bmacc_state = stateWAIT;
        pr_info("Device File Opened...!!!\n");
        return 0;
}

/*
** This function will be called when we close the Device file
*/
static int bm_release(struct inode *inode, struct file *file)
{
        bmacc_state = stateWAIT;
        pr_info("Device File Closed...!!!\n");
        return 0;
}

/*
** This function will be called when we read the Device file
*/
static ssize_t bm_read(struct file *filp, char __user *buf, size_t len, loff_t *off)
{
        struct work_data *writework;

        wait_event_interruptible(wait_queue_bm, wait_queue_flag != 0);

        switch (wait_queue_flag)
        {
        case 1:
                switch (bmacc_state)
                {
                case stateHSRECV:
                        if (copy_to_user(buf, &hsmask, 1))
                        {
                                pr_err("Data Read : Err!\n");
                        }
                        //pr_info("Sent MASK");
                        bmacc_state = stateMASKSENT;
                        wait_queue_flag = 0;
                        return 1;
                        break;
                case stateACK:
                        *replies = hs & hsmask;
                        if (copy_to_user(buf, replies, 1))
                        {
                                pr_err("Data Read : Err!\n");
                        }
                        //pr_info("Sent ACK");
                        bmacc_state = stateCONNECT;
                        wait_queue_flag = 0;
                        writework = kmalloc(sizeof(struct work_data), GFP_KERNEL);
                        INIT_WORK(&writework->work, work_handler);
                        writework->cs = 4;
                        queue_work(wq, &writework->work);
                        return 1;
                        break;
                }
                break;
        case 2:
                *replies = cmdKEEP;
                if (copy_to_user(buf, replies, 1))
                {
                        pr_err("Data Read : Err!\n");
                }
                wait_queue_flag = 0;
                return 1;
                break;
        case 3:
        case 4:
                writel((readl(bmptr + BMSTATES_O) | 0x1) & ~0x2, bmptr + BMSTATES_O); // Exec = 1 , Done = 0
                //pr_info("States: 0x%08x\n", readl(bmptr + BMSTATES_O));

                seqlen = 0;

                changes = readl(bmptr + BMCHANGES_O);

                if ((changes & CH_I0_RECV) != 0x0)
                {
                        if ((readl(bmptr + BMPL2PS_O) & PL2PS_I0_RECV) != 0x0)
                        {
                                *replies = cmdDRECVH | I0_ID;
                                //pr_info("Sent RECVH %08x\n", I0_ID);
                        }
                        else
                        {
                                *replies = cmdDRECVL | I0_ID;
                                //pr_info("Sent RECVL %08x\n", I0_ID);
                        }

                        if (copy_to_user(buf + seqlen, replies, 1))
                        {
                                //pr_err("Data Read : Err!\n");
                        }

                        seqlen = seqlen + 1;
                }

                if ((changes & CH_O0) != 0x0)
                {
                        value = readl(bmptr + BMO0_O);
                        *replies = cmdNEWVAL | O0_ID;

                        if (copy_to_user(buf + seqlen, replies, 1))
                        {
                                pr_err("Data Read : Err!\n");
                        }

                        seqlen = seqlen + 1;

                        if (copy_to_user(buf + seqlen, &value, REGSIZE))
                        {
                                pr_err("Data Read : Err!\n");
                        }
                        //pr_info("Sent NEWVAL %08x = %08x\n", cmdNEWVAL | O0_ID, value);

                        seqlen = seqlen + REGSIZE;
                }

                if ((changes & CH_O0_VALID) != 0x0)
                {
                        if ((readl(bmptr + BMPL2PS_O) & PL2PS_O0_VALID) != 0x0)
                        {
                                *replies = cmdDVALIDH | O0_ID;
                                //pr_info("Sent VALIDH %08x\n", O0_ID);
                        }
                        else
                        {
                                *replies = cmdDVALIDL | O0_ID;
                                //pr_info("Sent VALIDL %08x\n", O0_ID);
                        }

                        if (copy_to_user(buf + seqlen, replies, 1))
                        {
                                pr_err("Data Read : Err!\n");
                        }

                        seqlen = seqlen + 1;
                }

                writel(readl(bmptr + BMSTATES_O) | 0x2, bmptr + BMSTATES_O); // Done = 1
                //pr_info("States: 0x%08x\n", readl(bmptr + BMSTATES_O));

                //                writel(0x0, bmptr + BMSTATES_O);
                writel(readl(bmptr + BMSTATES_O) & ~0x1, bmptr + BMSTATES_O); // Exec = 0
                //pr_info("States: 0x%08x\n", readl(bmptr + BMSTATES_O));

                wait_queue_flag = 0;
                return seqlen;
                break;
        }
        pr_info("Data Read : Done!\n");
        return 0;
}

/*
** This function will be called when we write the Device file
*/
static ssize_t bm_write(struct file *filp, const char __user *buf, size_t len, loff_t *off)
{
        struct work_data *writework;

        if (copy_from_user(write_buffer, buf, len))
        {
                pr_err("data write error\n");
        }
        else
        {
                for (i = 0; i < len; i++)
                {
                        switch (bmacc_state)
                        {
                        case stateWAIT:
                                switch (write_buffer[i] & cmdMASK)
                                {
                                case cmdHANDSH:
                                        //pr_info("HandShake");
                                        hs = write_buffer[i];
                                        bmacc_state = stateHSRECV;
                                        wait_queue_flag = 0;
                                        writework = kmalloc(sizeof(struct work_data), GFP_KERNEL);
                                        INIT_WORK(&writework->work, work_handler);
                                        writework->cs = 1;
                                        queue_work(wq, &writework->work);
                                        break;
                                }
                                break;
                        case stateMASKSENT:
                                if (write_buffer[i] == (hs & hsmask))
                                {
                                        bmacc_state = stateACK;
                                        wait_queue_flag = 0;
                                        writework = kmalloc(sizeof(struct work_data), GFP_KERNEL);
                                        INIT_WORK(&writework->work, work_handler);
                                        writework->cs = 1;
                                        queue_work(wq, &writework->work);
                                }
                                else
                                {
                                        bmacc_state = stateWAIT;
                                }
                                break;
                        case stateCONNECT:
                                if (waitvalue > 0)
                                {
                                        *(value_buffer + waitvalue - 1) = write_buffer[i];
                                        waitvalue = waitvalue - 1;
                                        if (waitvalue == 0)
                                        {
                                                switch (regvalue)
                                                {
                                                case I0_ID:
                                                        writel(*value_buffer, bmptr + BMI0_O);
                                                        //pr_info("Received I%d - Value %d", 0, *value_buffer);
                                                        break;
                                                }
                                        }
                                }
                                else
                                {
                                        switch (write_buffer[i] & cmdMASK)
                                        {
                                        case cmdNEWVAL:
                                                switch (write_buffer[i] & ~cmdMASK)
                                                {
                                                case I0_ID:
                                                        regvalue = I0_ID;
                                                        waitvalue = REGSIZE;
                                                        break;
                                                }
                                                break;
                                        case cmdDVALIDH:
                                                switch (write_buffer[i] & ~cmdMASK)
                                                {
                                                case I0_ID:
                                                        //pr_info("Received I%d - Valid H", 0);
                                                        writel(readl(bmptr + BMPS2PL_O) | PS2PL_I0_VALID, bmptr + BMPS2PL_O);
                                                        break;
                                                }
                                                break;
                                        case cmdDVALIDL:
                                                switch (write_buffer[i] & ~cmdMASK)
                                                {
                                                case I0_ID:
                                                        //pr_info("Received I%d - Valid L", 0);
                                                        writel(readl(bmptr + BMPS2PL_O) & ~PS2PL_I0_VALID, bmptr + BMPS2PL_O);
                                                        break;
                                                }
                                                break;
                                        case cmdDRECVH:
                                                switch (write_buffer[i] & ~cmdMASK)
                                                {
                                                case O0_ID:
                                                        //pr_info("Received O%d - Recv H", 0);
                                                        writel(readl(bmptr + BMPS2PL_O) | PS2PL_O0_RECV, bmptr + BMPS2PL_O);
                                                        break;
                                                }
                                                break;
                                        case cmdDRECVL:
                                                switch (write_buffer[i] & ~cmdMASK)
                                                {
                                                case O0_ID:
                                                        //pr_info("Received O%d - Recv L", 0);
                                                        writel(readl(bmptr + BMPS2PL_O) & ~PS2PL_O0_RECV, bmptr + BMPS2PL_O);
                                                        break;
                                                }
                                                break;
                                        case cmdKEEP:
                                                // pr_info("KEEP received")
                                                writework = kmalloc(sizeof(struct work_data), GFP_KERNEL);
                                                INIT_WORK(&writework->work, work_handler);
                                                writework->cs = 2;
                                                queue_work(wq, &writework->work);
                                                break;
                                        }
                                }
                                break;
                        }
                }
        }

        return len;
}

/*
** Module Init function
*/
static int __init bmmod_init(void)
{

        wq = create_workqueue("bm_queue");

        bmptr = ioremap(BMBASEADDR, BMLEN);

        // pr_info("%x", readl(bmptr));
        // pr_info("%x", readl(bmptr + 4));
        // pr_info("%x", readl(bmptr + 8));
        // pr_info("%x", readl(bmptr + 12));
        // pr_info("%x", readl(bmptr + 16));
        // pr_info("%x", readl(bmptr + 20));

        /*Allocating Major number for the BM device */
        if ((alloc_chrdev_region(&dev, 0, 1, "bm_dev")) < 0)
        {
                pr_info("Cannot allocate major number\n");
                return -1;
        }
        pr_info("Major = %d Minor = %d \n", MAJOR(dev), MINOR(dev));

        /*Creating cdev structure*/
        cdev_init(&bm_cdev, &fops);

        /*Adding character device to the system*/
        if ((cdev_add(&bm_cdev, dev, 1)) < 0)
        {
                pr_info("Cannot add the device to the system\n");
                goto r_class;
        }

        /*Creating struct class*/
        if ((dev_class = class_create(THIS_MODULE, "bm_class")) == NULL)
        {
                pr_info("Cannot create the struct class\n");
                goto r_class;
        }

        /*Creating device*/
        if ((device_create(dev_class, NULL, dev, NULL, "bm")) == NULL)
        {
                pr_info("Cannot create the Device 1\n");
                goto r_device;
        }

        // Initialize wait queue
        init_waitqueue_head(&wait_queue_bm);

        /*Creating Physical memory*/
        if ((write_buffer = kmalloc(write_buff_size, GFP_KERNEL)) == 0)
        {
                pr_info("Cannot allocate memory in kernel\n");
                goto r_device;
        }

        if ((value_buffer = kmalloc(4, GFP_KERNEL)) == 0)
        {
                pr_info("Cannot allocate memory in kernel\n");
                goto r_device;
        }

        if ((replies = kmalloc(bm_reg_size + 1, GFP_KERNEL)) == 0)
        {
                pr_info("Cannot allocate memory in kernel\n");
                goto r_device;
        }

        if (request_irq(IRQ_NO, irq_handler, IRQF_SHARED, "bm_irq", (void *)(irq_handler)))
        {
                pr_err("my_device: cannot register IRQ ");
                goto irq;
        }
        if (irq_set_irq_type(IRQ_NO, IRQ_TYPE_EDGE_RISING))
        {
                goto irq;
        }

        pr_info("Device Driver Insert...Done!!!\n");

        return 0;

irq:
        free_irq(IRQ_NO, (void *)(irq_handler));
r_device:
        class_destroy(dev_class);
r_class:
        unregister_chrdev_region(dev, 1);
        return -1;
}

/*
** Module exit function
*/
static void __exit bmmod_exit(void)
{
        kfree(value_buffer);
        kfree(write_buffer);
        kfree(replies);
        wait_queue_flag = 5;
        wake_up_interruptible(&wait_queue_bm);
        device_destroy(dev_class, dev);
        class_destroy(dev_class);
        cdev_del(&bm_cdev);
        unregister_chrdev_region(dev, 1);
        free_irq(IRQ_NO, (void *)(irq_handler));
        iounmap(bmptr);
        flush_workqueue(wq);
        destroy_workqueue(wq);
        pr_info("Device Driver Remove...Done!!!\n");
}

module_init(bmmod_init);
module_exit(bmmod_exit);

MODULE_LICENSE("GPL");
MODULE_AUTHOR("BondMachine <bondmachine@fisica.unipg.it>");
MODULE_DESCRIPTION("");
MODULE_VERSION("1.0");
`
)
