import logging


class DeliveryBoy:
    def __init__(self, env, id, dbManager, x, y, speed):
        self.id = id
        self.speed = speed
        self.name = "DB-%d" % id
        self.x = x
        self.y = y
        self.env = env
        self.dbManager = dbManager
        self.lastDeliveryTime = 0
        logging.debug("%s: Setup at X:%d,Y:%d" % (self.name, x, y))

    def deliver(self, order):
        logging.debug("%s (O%d): received at %d" % (self.name, order.id, self.env.now))
        yield self.env.process(self.drive_to_restaurant(order))
        yield self.env.process(self.pickup_food(order))
        yield self.env.process(self.drive_to_customer(order))

    def drive_to_restaurant(self, order):
        distance_to_restaurant = order.distance_to_restaurant(self.x, self.y)
        time_required = self.time_required(distance_to_restaurant)
        logging.info(
            "%s (O%d):STARTING_PICKUP Distance: %d TimeRequired: %d at %d" % (
                self.name, order.id, distance_to_restaurant, time_required, self.env.now))
        yield self.env.timeout(time_required)
        logging.debug("%s (O%d): Reached Restaurant at %d" % (self.name, order.id, self.env.now))
        self.x, self.y = order.restaurant.x, order.restaurant.y

    def pickup_food(self, order):
        yield self.env.process(order.restaurant.handover_food(order))
        logging.debug("%s (O%d): PICKED_FOOD at %d" % (self.name, order.id, self.env.now))

    def drive_to_customer(self, order):
        distance_to_customer = order.distance_to_customer(self.x, self.y)
        time_required = self.time_required(distance_to_customer)
        yield self.env.timeout(time_required)
        logging.info(
            "%s (O%d): #REACHED_CUSTOMER# Distance: %d TimeTaken: %d at %d" % (
                self.name, order.id, distance_to_customer, time_required, self.env.now))
        self.lastDeliveryTime = self.env.now
        self.dbManager.reportOrderServed(self, self.env.now - order.orderTime)

    def time_required(self, distance):
        # distance = speed x time
        return distance / self.speed

    def getCost(self, order, RESTAURANT_WEIGHT, IDLE_WEIGHT):
        time_to_restaurant = self.time_required(order.distance_to_restaurant(self.x, self.y))
        idleTime = self.env.now - self.lastDeliveryTime

        # Weighted Average Cost based on time req to reach restaurant and idle time
        cost = ((time_to_restaurant * RESTAURANT_WEIGHT) + (idleTime * IDLE_WEIGHT)) / (
                RESTAURANT_WEIGHT + IDLE_WEIGHT)
        return cost
