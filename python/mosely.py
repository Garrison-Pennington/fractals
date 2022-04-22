import numpy as np
import matplotlib.pyplot as plt
import matplotlib

from python import Fractal


class Mosely(Fractal):
    def __init__(self, steps=3, light=False):
        self.fractal = np.array([[[1]]])
        self.steps = steps
        self.sub = np.array([
            [[0, 1, 0], [1, 1, 1], [0, 1, 0]],
            [[1, 1, 1], [0 if light else 1 for _ in range(3)], [1, 1, 1]],
            [[0, 1, 0], [1, 1, 1], [0, 1, 0]],
        ])

    def step(self, fractal):
        h, r, d = fractal.shape
        new_cube = np.zeros((h*3, r*3, d*3))
        for h, r, d in zip(*fractal.nonzero()):
            new_cube[3*h:3*h+3, 3*r:3*r+3, 3*d:3*d+3] = self.sub
        return new_cube

    def render(self):
        self.fractal, _ = self.n_steps(self.steps)
        fig = plt.figure()
        ax = fig.add_subplot(projection='3d')
        ax.voxels(self.fractal)
        plt.show()


matplotlib.use("TkAgg")

f = Mosely()
f.render()
