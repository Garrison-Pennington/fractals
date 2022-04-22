import random
import string

import numpy as np
import matplotlib.pyplot as plt
import matplotlib
from tqdm import tqdm

from python import Fractal

matplotlib.use("TkAgg")


class Mandelbrot(Fractal):
    def __init__(self, steps=100):
        self.c = None
        self.steps = steps
        self.fractal = 0

    def step(self, fractal):
        return fractal * fractal + self.c

    def exit_condition(self, fractal):
        return abs(fractal) > 2

    def score(self, fractal, step):
        return step / self.steps

    def render(self, image_shape, real_bounds, imaginary_bounds, steps):
        img = np.zeros(image_shape)
        h, w = image_shape
        real = np.linspace(*real_bounds, w)
        imaginary = np.linspace(*imaginary_bounds, h)
        for ii, i in tqdm(list(enumerate(imaginary))):
            for ir, r in enumerate(real):
                self.c = complex(r, i)
                self.fractal, step = self.n_steps(steps)
                img[ii, ir] = self.score(self.fractal, step)
                self.fractal = 0
        return img


m = Mandelbrot()

img = m.render((100000, 100000), (-2, 2), (-2, 2), 100)
m.save(img, "mandelbrot_"
       + "".join(random.choices(string.ascii_lowercase, k=16))+".png")
