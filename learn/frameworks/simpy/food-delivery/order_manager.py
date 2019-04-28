import logging
import random

from models.order import Order, Customer


class OrderManager:
    def __init__(self, env, dbManager, resManager, xy_generator, orderConfig):
        self.env = env
        self.genId = 0
        self.orderConfig = orderConfig
        self.dbManager = dbManager
        self.resManager = resManager
        self.xy_generator = xy_generator

    def place_order(self, order):
        # Tell Restaurant to Prepare Food
        self.env.process(order.restaurant.prepare_food(order))

        # Tell Delivery Boy to Deliver Food
        self.env.process(self.dbManager.deliverOrder(order))

    def order_generator(self):
        interval = self.orderConfig["generateInterval"]
        burstConfig = self.orderConfig["burst"]
        logging.critical("-------- Starting Order Generation ----------")
        while True:
            yield self.env.timeout(random.randint(interval - 2, interval + 2))
            self.genId += 1

            dish = self.resManager.get_random_dish()
            restaurant = self.resManager.get_random_restaurant()

            if random.randint(0, 100) < 100 - burstConfig["percent"]:
                # Steady Flow of Orders 90% of the time
                self.generate_order(self.genId, dish, restaurant)
            else:
                # Generate Order Bursts for Same Restaurant
                orderCount = random.randint(burstConfig["min"], burstConfig["max"])
                logging.info("Generating Order Burst: %d Orders" % orderCount)
                for i in range(orderCount):
                    self.generate_order(self.genId, dish, restaurant)
                    self.genId += 1

    def generate_order(self, id, dish, restaurant):
        x, y = self.xy_generator.next()
        customer = Customer(x, y)

        logging.info("NEW_ORDER (O%d): Dish %d Restaurant %d Customer (%d,%d) at %d" % (
            id, dish.id, restaurant.id, x, y, self.env.now))
        self.place_order(Order(id, restaurant, dish, customer, self.env.now))
