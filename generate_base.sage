# This is as sketch of how to decaffeinate Curve448 based on Curve25519
# check here: https://cr.yp.to/highspeed/naclcrypto-20090310.pdf

F = GF(2^255-19)
# F = GF(2^448 - 2^224 - 1)
# d = -39081
d = -121665
# constructs an ec from Weitrass a-coefficients
# sage: E = EllipticCurve([1,2,3,4,5]); E
# Elliptic Curve defined by y^2 + x*y + 3*y = x^3 + 2*x^2 + 4*x + 5 over Rational Field
# a * x2 + y2 = 1 + d * x2  y2
# how? maybe with the ring
M = EllipticCurve(F,[0,2-4*d,0,1,0])

# check
sqrtN1 = sqrt(F(-1))

debugging = True
def debug_print(foo):
    if debugging: print foo

def maybe(): return randint(0,1)

# check
def qpositive(x):
    return int(x) <= (2^255-19-1)/2

def M_to_X(x, y):
    s = sqrt(x)
    if s == 0: t = 1
    else: t = y/s

    X,Y = 2*s / (1+s^2), (1-s^2) / t # This is phi_a(s, t) page 7
    if maybe(): X,Y = -X,-Y
    if maybe(): X,Y = Y,-X
    # OK, have point in ed
    return X,Y

def M_to_E(P):
    # P must be even
    (x,y) = P.xy()
    assert x.is_square()

    s = sqrt(x)
    if s == 0: t = 1
    else: t = y/s

    X = 2*s / (1+s^2)
    Y = (1-(s^2)) / t # This is phi_a(s, t) page 7
    if maybe(): X,Y = -X,-Y
    if maybe(): X,Y = Y,-X
    # OK, have point in ed
    return X,Y

def decaf_encode_from_E(X,Y):
    assert X^2 + Y^2 == 1 + d*X^2*Y^2 # curve
    if not qpositive(X*Y): X,Y = Y,-X
    assert qpositive(X*Y)

    assert (1-X^2).is_square()
    sx = sqrt(1-X^2)
    tos = -2*sx/X/Y
    if not qpositive(tos): sx = -sx
    s = (1 + sx) / X
    if not qpositive(s): s = -s

    return s

def isqrt(x):
    ops = [(1,2),(1,2),(3,1),(6,0),(1,2),(12,1),(25,1),(25,1),(50,0),(125,0),(2,2),(1,2)]
    st = [x,x,x]
    for i,(sh,add) in enumerate(ops):
        od = i&1
        st[od] = st[od^^1]^(2^sh)*st[add]
    # assert st[2] == x^(2^252-3)

    assert st[1] == 1 or st[1] == -1
    if st[1] == 1: return st[0]
    else: return st[0] * sqrtN1

def decaf_encode_from_E_c(X,Y):
    Z = F.random_element()
    T = X*Y*Z
    X = X*Z
    Y = Y*Z
    assert X^2 + Y^2 == Z^2 + d*T^2

    # Precompute
    sd = sqrt(F(1-d))

    zx = Z^2-X^2
    TZ = T*Z
    assert zx.is_square
    ooAll = isqrt(zx*TZ^2)
    osx = ooAll * TZ
    ooTZ = ooAll * zx * osx

    floop = qpositive(T^2 * ooTZ)
    if floop:
        frob = zx * ooTZ
    else:
        frob = sd
        Y = -X

    osx *= frob

    if qpositive(-2*osx*Z) != floop: osx = -osx
    s = Y*(ooTZ*Z + osx)
    if not qpositive(s): s = -s

    return s

def is_rotation((X,Y),(x,y)):
    return x*Y == X*y or x*X == -y*Y

def decaf_decode_to_E(s):
    assert qpositive(s)
    t = sqrt(s^4 + (2-4*d)*s^2 + 1)
    if not qpositive(t/s): t = -t
    X,Y = 2*s / (1+s^2), (1-s^2) / t
    assert qpositive(X*Y)
    return X,Y

def decaf_decode_to_E_c(s):
    assert qpositive(s)

    s2 = s^2
    s21 = 1+s2
    t2 = s21^2 - 4*d*s2

    alt  = s21*s
    the  = isqrt(t2*alt^2)
    oot  = the * alt
    the *= t2
    tos  = the * s21
    X = 2 * (tos-the) * oot
    Y = (1-s2) * oot

    if not qpositive(tos): Y = -Y
    assert qpositive(X*Y)

    return X,Y

def test():
    P = 2*M.random_point()
    X,Y = M_to_E(P)
    s = decaf_encode_from_E(X,Y)
    assert s == decaf_encode_from_E_c(X,Y)
    XX,YY = decaf_decode_to_E(s)
    XX2,YY2 = decaf_decode_to_E_c(s)
    assert is_rotation((X,Y),(XX,YY))
    assert is_rotation((X,Y),(XX2,YY2))

P = 2*M.random_point()
#X,Y = M_to_E(P)
#print("the X", X)
#print("the Y", Y)
#print decaf_encode_from_E(X, Y)
#print M.base_ring()
#print M.weierstrass_p(prec=10)
x = 117812161263436946737282484343310064665180535357016373416879082147939404277809514858788439644911793978499419995990477371552926308078495
y = 19
print(M_to_X(x, y))
# print("the x", X1)
# print("the y", Y1)
