import simpy

from entities.delivery_boy import DeliveryBoy


class DeliveryBoyManager:
    def __init__(self, env, count):
        self.env = env
        self.freePool = simpy.Store(env, count)
        self.dbGenerator = self.delivery_boy_generator(0, self.freePool)
        self.orderServed = 0
        for i in range(count):
            self.freePool.put(self.dbGenerator.next())

    def deliverOrder(self, order):
        boy = yield self.freePool.get()
        self.env.process(boy.deliver(order))
        self.orderServed += 1

    def delivery_boy_generator(self, i, pool):
        while True:
            i += 1
            yield DeliveryBoy(self.env, i, pool)

    def printOrdersServed(self):
        print "Orders Served: %d" % self.orderServed
