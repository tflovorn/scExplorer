from mpl_toolkits.mplot3d import Axes3D
from matplotlib import cm
from matplotlib.ticker import LinearLocator, FormatStrFormatter
import matplotlib.pyplot as plt
import numpy as np

def plotSymFactor3D(alpha):
    fig = plt.figure()
    ax = fig.gca(projection='3d')
    X = np.arange(-np.pi, np.pi+0.001, 0.1)
    Y = np.arange(-np.pi, np.pi+0.001, 0.1)
    X, Y = np.meshgrid(X, Y)
    Z = np.absolute(np.sin(X) + float(alpha)*np.sin(Y))
    surf = ax.plot_surface(X, Y, Z, rstride=1, cstride=1, cmap=cm.coolwarm,
            linewidth=0, antialiased=False)
    ax.set_zlim(0.0, 2.0)

    ax.zaxis.set_major_locator(LinearLocator(10))
    ax.zaxis.set_major_formatter(FormatStrFormatter("%.02f"))

    plt.show()

def plotSymFactor(alpha, fname, largeFont):
    X = np.arange(-np.pi, np.pi+0.001, 0.01)
    Y = np.arange(-np.pi, np.pi+0.001, 0.01)
    X, Y = np.meshgrid(X, Y)
    Z = np.absolute(np.sin(X) + float(alpha)*np.sin(Y))

    fig = plt.figure()

    if largeFont:
        plt.xlabel("$k_x$", fontsize=20)
        plt.ylabel("$k_y$", fontsize=20)
        plt.xticks(fontsize=20)
        plt.yticks(fontsize=20)
    else:
        plt.xlabel("$k_x$")
        plt.ylabel("$k_y$")

    N = 30
    p = plt.contourf(X, Y, Z, N, cmap=cm.binary)

    cb = fig.colorbar(p)
    
    if largeFont:
        cb.ax.tick_params(labelsize=20)

    plt.savefig(fname + ".png", bbox_inches="tight", dpi=200)


if __name__ == "__main__":
    plotSymFactor(-1, "sym_alpha_m1", False)
    plotSymFactor(1, "sym_alpha_p1", False)
