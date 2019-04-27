import logging


class DeliveryBoy:
    def __init__(self, env, id, pool, x, y, speed):
        self.id = id
        self.speed = speed
        self.name = "DB-%d" % id
        self.x = x
        self.y = y
        self.env = env
        self.pool = pool
        logging.debug("%s: Setup at X:%d,Y:%d" % (self.name, x, y))

    def deliver(self, order):
        logging.debug("%s (O%d): received at %d" % (self.name, order.id, self.env.now))
        yield self.env.process(self.drive_to_restaurant(order))
        yield self.env.process(self.pickup_food(order))
        yield self.env.process(self.drive_to_customer(order))
        yield self.env.process(self.handover_food(order))

    def drive_to_restaurant(self, order):
        distance_to_restaurant = order.distance_to_restaurant(self.x, self.y)
        logging.info(
            "%s (O%d):STARTING_PICKUP Distance: %d TimeRequired: %d at %d" % (
                self.name, order.id, distance_to_restaurant, self.time_required(distance_to_restaurant), self.env.now))
        yield self.env.timeout(2)
        logging.debug("%s (O%d): Reached Restaurant at %d" % (self.name, order.id, self.env.now))
        self.x, self.y = order.restaurant.x, order.restaurant.y

    def pickup_food(self, order):
        yield self.env.process(order.restaurant.handover_food(order))
        logging.debug("%s (O%d): PICKED_FOOD at %d" % (self.name, order.id, self.env.now))

    def drive_to_customer(self, order):
        yield self.env.timeout(order.customer_drive_time())
        logging.debug("%s (O%d): REACHED_CUSTOMER at %d" % (self.name, order.id, self.env.now))

    def handover_food(self, order):
        yield self.env.timeout(order.customer_handover_time())
        logging.info("%s (O%d): #CUSTOMER_HANDOVER_DONE#  at %d" % (self.name, order.id, self.env.now))
        yield self.pool.put(self)

    def time_required(self, distance):
        # distance = speed x time
        return distance / self.speed
