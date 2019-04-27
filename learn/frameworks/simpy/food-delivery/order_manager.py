class OrderManager:
    def __init__(self, env, dbManager):
        self.env = env
        self.dbManager = dbManager

    def place_order(self, order):
        # Tell Restaurant to Prepare Food
        self.env.process(order.restaurant.prepare_food(order))

        # Tell Delivery Boy to Deliver Food
        boy = self.dbManager.getdeliveryboy()
        self.env.process(boy.deliver(order))
