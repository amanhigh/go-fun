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


class Customer:
    def __init__(self, x, y):
        self.x = x
        self.y = y


class Order:
    def __init__(self, id, restaurant, dish, customer, orderTime):
        self.id = id
        self.restaurant = restaurant
        self.dish = dish
        self.customer = customer
        self.orderTime = orderTime

    def distance_to_restaurant(self, x, y):
        return math.hypot(x - self.restaurant.x, y - self.restaurant.y)

    def distance_to_customer(self, x, y):
        return math.hypot(x - self.customer.x, y - self.customer.y)
