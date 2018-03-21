# Run using docker run -it -v `pwd`/large_models:/models mxnet/python python /models/alexnet.py

from __future__ import print_function
import mxnet as mx
from mxnet import nd, autograd
from mxnet import gluon
import numpy as np
from mxnet import ndarray as F
from mxnet import sym

mx.random.seed(1)
ctx = mx.cpu()

mult = 4

alex_net = gluon.nn.Sequential()
with alex_net.name_scope():
    #  First convolutional layer
    alex_net.add(gluon.nn.Conv2D(channels=mult*96, kernel_size=mult*11, strides=(4,4), activation='relu'))
    alex_net.add(gluon.nn.MaxPool2D(pool_size=3, strides=2))
    #  Second convolutional layer
    alex_net.add(gluon.nn.Conv2D(channels=mult*192, kernel_size=mult*5, activation='relu'))
    alex_net.add(gluon.nn.MaxPool2D(pool_size=3, strides=(2,2)))
    # Third convolutional layer
    alex_net.add(gluon.nn.Conv2D(channels=mult*384, kernel_size=mult*3, activation='relu'))
    # Fourth convolutional layer
    alex_net.add(gluon.nn.Conv2D(channels=mult*384, kernel_size=mult*3, activation='relu'))
    # Fifth convolutional layer
    alex_net.add(gluon.nn.Conv2D(channels=mult*256, kernel_size=3, activation='relu'))
    alex_net.add(gluon.nn.MaxPool2D(pool_size=3, strides=2))
    # Flatten and apply fullly connected layers
    alex_net.add(gluon.nn.Flatten())
    alex_net.add(gluon.nn.Dense(4096, activation="relu"))
    alex_net.add(gluon.nn.Dense(4096, activation="relu"))
    alex_net.add(gluon.nn.Dense(10))

alex_net.collect_params().initialize(mx.init.Xavier(magnitude=2.24), ctx=ctx, force_reinit=True)

# alex_net.initialize(ctx=ctx, force_reinit=True)

alex_net(mx.nd.random_normal(shape=(16, 3, mult*224, mult*224)))

with open('/models/alexnet_%d.json' % mult, 'w') as file:
  data = sym.var('data')
  sy = alex_net(data)
  file.write(sy.tojson())

alex_net.save_params('/models/alexnet_%d.params' % mult)

