class OrderManager:
    def __init__(self, env, deliveryBoy):
        self.env = env
        self.deliveryBoy = deliveryBoy

    def place_order(self, order):
        # Tell Restaurant to Prepare Food
        # self.env.process(order.restaurant.prepare_food(order))

        # Tell Delivery Boy to Deliver Food
        self.env.process(self.deliveryBoy.deliver(order))
