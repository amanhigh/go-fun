import logging
import math


class Dish:
    def __init__(self, id, time):
        self.id = id
        self.time = time
        logging.debug("Dish %d: Cooktime-%d" % (id, time))

    # Assumption Prep time for dish
    # is same across all restaurants.
    def prep_time(self):
        return self.time


class Order:
    def __init__(self, id, restaurant, dish):
        self.id = id
        self.restaurant = restaurant
        self.dish = dish

    def distance_to_restaurant(self, x, y):
        hypot = math.hypot(x - self.restaurant.x, y - self.restaurant.y)
        return hypot

    def customer_drive_time(self):
        return 2

    def customer_handover_time(self):
        return 2
