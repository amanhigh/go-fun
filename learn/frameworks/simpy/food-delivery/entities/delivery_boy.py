import logging


class DeliveryBoy:
    def __init__(self, env, id, pool):
        self.id = id
        self.name = "DB-%d" % id
        self.env = env
        self.pool = pool

    def deliver(self, order):
        logging.debug("%s (O%d): received at %d" % (self.name, order.id, self.env.now))
        yield self.env.process(self.drive_to_restaurant(order))
        yield self.env.process(self.pickup_food(order))
        yield self.env.process(self.drive_to_customer(order))
        yield self.env.process(self.handover_food(order))

    def drive_to_restaurant(self, order):
        logging.info("%s (O%d): Starting Pickup at %d" % (self.name, order.id, self.env.now))
        yield self.env.timeout(2)
        logging.debug("%s (O%d): Reached Restaurant at %d" % (self.name, order.id, self.env.now))

    def pickup_food(self, order):
        yield self.env.process(order.restaurant.handover_food(order))
        logging.debug("%s (O%d): Picked Food at %d" % (self.name, order.id, self.env.now))

    def drive_to_customer(self, order):
        yield self.env.timeout(order.customer_drive_time())
        logging.debug("%s (O%d): Reached Customer at %d" % (self.name, order.id, self.env.now))

    def handover_food(self, order):
        yield self.env.timeout(order.customer_handover_time())
        logging.info("%s (O%d): #Customer Handover done#  at %d" % (self.name, order.id, self.env.now))
        yield self.pool.put(self)
