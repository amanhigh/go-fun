class DeliveryBoy:
    def __init__(self, env):
        self.env = env

    def deliver(self, order):
        print ("Got Order: %d at %d" % (order.id, self.env.now))
        yield self.env.timeout(2)
