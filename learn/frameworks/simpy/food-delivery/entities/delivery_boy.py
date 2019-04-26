import logging


class DeliveryBoy:
    def __init__(self, env, id):
        self.id = id
        self.name = "DB-%d" % id
        self.env = env

    def deliver(self, order):
        logging.info("%s: New Delivery Order: %d at %d" % (self.name, order.id, self.env.now))
        yield self.env.process(self.drive_to_restaurant(order))
        yield self.env.process(self.pickup_food(order))
        yield self.env.process(self.drive_to_customer(order))
        yield self.env.process(self.handover_food(order))

    def drive_to_restaurant(self, order):
        yield self.env.timeout(2)
        logging.info("%s: Reached Restaurant Order: %d at %d" % (self.name, order.id, self.env.now))

    def pickup_food(self, order):
        yield self.env.process(order.restaurant.handover_food(order))
        logging.info("%s: Picked Food Order: %d at %d" % (self.name, order.id, self.env.now))

    def drive_to_customer(self, order):
        yield self.env.timeout(2)
        logging.info("%s: Reached Customer Order: %d at %d" % (self.name, order.id, self.env.now))

    def handover_food(self, order):
        yield self.env.timeout(2)
        logging.info("%s: Handedover Food Order: %d at %d" % (self.name, order.id, self.env.now))
