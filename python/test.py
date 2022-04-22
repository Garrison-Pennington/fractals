import math
from PIL import Image

import numpy as np
import matplotlib.pyplot as plt
import matplotlib
from tqdm import tqdm

from python import Fractal


class Test1(Fractal):
    def __init__(self, steps=100):
        self.c = None
        self.steps = steps
        self.fractal = 0

    def step(self, fractal):
        if abs(fractal) <= 2:
            fractal *= self.c
        else:
            fractal -= self.c
        return fractal

    def exit_condition(self, fractal):
        return abs(fractal) <= abs(2 - self.c/2)

    def score(self, fractal, step):
        return step / self.steps

    def render(self, image_shape, real_bounds, imaginary_bounds, steps=100):
        img = np.zeros(image_shape)
        real = np.linspace(*real_bounds, image_shape[1])
        imaginary = np.linspace(*imaginary_bounds, image_shape[0])
        for ir, r in tqdm(list(enumerate(real))):
            for ii, i in enumerate(imaginary):
                self.c = complex(r, i)
                self.fractal = 2
                z, n = self.n_steps(steps)
                img[ii, ir] = self.score(z, n)
        return img


class ComplexFib(Fractal):
    def __init__(self, steps=100):
        self.fractal = (1, 1)
        self.steps = 100

    def step(self, fractal):
        last, cur = fractal
        sub = last + cur
        return cur, sub * sub

    def exit_condition(self, fractal):
        last, cur = fractal
        return abs(cur) >= 2

    def score(self, n):
        return n / self.steps

    def render(self, image_shape, real_bounds, imaginary_bounds):
        h, w = image_shape
        img = np.zeros(image_shape)
        real = np.linspace(*real_bounds, w)
        imaginary = np.linspace(*imaginary_bounds, h)
        for ir, r in tqdm(list(enumerate(real))):
            for ii, i in enumerate(imaginary):
                c = complex(r, i)
                self.fractal = (c, c)
                val, n = self.n_steps(self.steps)
                img[ii, ir] = self.score(n)
        return img


f = ComplexFib()
img = f.render((10000, 10000), (-.5, .5), (-.5, .5))
f.save(img, "ComplexFib.png")
