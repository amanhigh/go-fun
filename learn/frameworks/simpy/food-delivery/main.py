import simpy
from entities.delivery_boy import DeliveryBoy
from models.order import Order

env = simpy.Environment()
boy = DeliveryBoy(env)
order = Order(1)

# Driver Interrupts last Car
env.process(boy.deliver(order))
env.run(until=20)
