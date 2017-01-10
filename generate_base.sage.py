
# This file was *autogenerated* from the file generate_base.sage
from sage.all_cmdline import *   # import sage library

_sage_const_3 = Integer(3); _sage_const_2 = Integer(2); _sage_const_1 = Integer(1); _sage_const_0 = Integer(0); _sage_const_6 = Integer(6); _sage_const_4 = Integer(4); _sage_const_255 = Integer(255); _sage_const_121665 = Integer(121665); _sage_const_125 = Integer(125); _sage_const_12 = Integer(12); _sage_const_10 = Integer(10); _sage_const_117812161263436946737282484343310064665180535357016373416879082147939404277809514858788439644911793978499419995990477371552926308078495 = Integer(117812161263436946737282484343310064665180535357016373416879082147939404277809514858788439644911793978499419995990477371552926308078495); _sage_const_19 = Integer(19); _sage_const_50 = Integer(50); _sage_const_25 = Integer(25)# This is as sketch of how to decaffeinate Curve448 based on Curve25519
# check here: https://cr.yp.to/highspeed/naclcrypto-20090310.pdf

F = GF(_sage_const_2 **_sage_const_255 -_sage_const_19 )
# F = GF(2^448 - 2^224 - 1)
# d = -39081
d = -_sage_const_121665 
# constructs an ec from Weitrass a-coefficients
# sage: E = EllipticCurve([1,2,3,4,5]); E
# Elliptic Curve defined by y^2 + x*y + 3*y = x^3 + 2*x^2 + 4*x + 5 over Rational Field
# a * x2 + y2 = 1 + d * x2  y2
# how? maybe with the ring
M = EllipticCurve(F,[_sage_const_0 ,_sage_const_2 -_sage_const_4 *d,_sage_const_0 ,_sage_const_1 ,_sage_const_0 ])

# check
sqrtN1 = sqrt(F(-_sage_const_1 ))

debugging = True
def debug_print(foo):
    if debugging: print foo

def maybe(): return randint(_sage_const_0 ,_sage_const_1 )

# check
def qpositive(x):
    return int(x) <= (_sage_const_2 **_sage_const_255 -_sage_const_19 -_sage_const_1 )/_sage_const_2 

def M_to_X(x, y):
    # P must be even
    s = sqrt(x)
    if s == _sage_const_0 : t = _sage_const_1 
    else: t = y/s

    X,Y = _sage_const_2 *s / (_sage_const_1 +s**_sage_const_2 ), (_sage_const_1 -s**_sage_const_2 ) / t # This is phi_a(s, t) page 7
    if maybe(): X,Y = -X,-Y
    if maybe(): X,Y = Y,-X
    # OK, have point in ed
    return X,Y

def M_to_E(P):
    # P must be even
    (x,y) = P.xy()
    assert x.is_square()

    s = sqrt(x)
    if s == _sage_const_0 : t = _sage_const_1 
    else: t = y/s

    X,Y = _sage_const_2 *s / (_sage_const_1 +s**_sage_const_2 ), (_sage_const_1 -s**_sage_const_2 ) / t # This is phi_a(s, t) page 7
    if maybe(): X,Y = -X,-Y
    if maybe(): X,Y = Y,-X
    # OK, have point in ed
    return X,Y

def decaf_encode_from_E(X,Y):
    assert X**_sage_const_2  + Y**_sage_const_2  == _sage_const_1  + d*X**_sage_const_2 *Y**_sage_const_2  # curve
    if not qpositive(X*Y): X,Y = Y,-X
    assert qpositive(X*Y)

    assert (_sage_const_1 -X**_sage_const_2 ).is_square()
    sx = sqrt(_sage_const_1 -X**_sage_const_2 )
    tos = -_sage_const_2 *sx/X/Y
    if not qpositive(tos): sx = -sx
    s = (_sage_const_1  + sx) / X
    if not qpositive(s): s = -s

    return s

def isqrt(x):
    ops = [(_sage_const_1 ,_sage_const_2 ),(_sage_const_1 ,_sage_const_2 ),(_sage_const_3 ,_sage_const_1 ),(_sage_const_6 ,_sage_const_0 ),(_sage_const_1 ,_sage_const_2 ),(_sage_const_12 ,_sage_const_1 ),(_sage_const_25 ,_sage_const_1 ),(_sage_const_25 ,_sage_const_1 ),(_sage_const_50 ,_sage_const_0 ),(_sage_const_125 ,_sage_const_0 ),(_sage_const_2 ,_sage_const_2 ),(_sage_const_1 ,_sage_const_2 )]
    st = [x,x,x]
    for i,(sh,add) in enumerate(ops):
        od = i&_sage_const_1 
        st[od] = st[od^_sage_const_1 ]**(_sage_const_2 **sh)*st[add]
    # assert st[2] == x^(2^252-3)

    assert st[_sage_const_1 ] == _sage_const_1  or st[_sage_const_1 ] == -_sage_const_1 
    if st[_sage_const_1 ] == _sage_const_1 : return st[_sage_const_0 ]
    else: return st[_sage_const_0 ] * sqrtN1

def decaf_encode_from_E_c(X,Y):
    Z = F.random_element()
    T = X*Y*Z
    X = X*Z
    Y = Y*Z
    assert X**_sage_const_2  + Y**_sage_const_2  == Z**_sage_const_2  + d*T**_sage_const_2 

    # Precompute
    sd = sqrt(F(_sage_const_1 -d))

    zx = Z**_sage_const_2 -X**_sage_const_2 
    TZ = T*Z
    assert zx.is_square
    ooAll = isqrt(zx*TZ**_sage_const_2 )
    osx = ooAll * TZ
    ooTZ = ooAll * zx * osx

    floop = qpositive(T**_sage_const_2  * ooTZ)
    if floop:
        frob = zx * ooTZ
    else:
        frob = sd
        Y = -X

    osx *= frob

    if qpositive(-_sage_const_2 *osx*Z) != floop: osx = -osx
    s = Y*(ooTZ*Z + osx)
    if not qpositive(s): s = -s

    return s

def is_rotation((X,Y),(x,y)):
    return x*Y == X*y or x*X == -y*Y

def decaf_decode_to_E(s):
    assert qpositive(s)
    t = sqrt(s**_sage_const_4  + (_sage_const_2 -_sage_const_4 *d)*s**_sage_const_2  + _sage_const_1 )
    if not qpositive(t/s): t = -t
    X,Y = _sage_const_2 *s / (_sage_const_1 +s**_sage_const_2 ), (_sage_const_1 -s**_sage_const_2 ) / t
    assert qpositive(X*Y)
    return X,Y

def decaf_decode_to_E_c(s):
    assert qpositive(s)

    s2 = s**_sage_const_2 
    s21 = _sage_const_1 +s2
    t2 = s21**_sage_const_2  - _sage_const_4 *d*s2

    alt  = s21*s
    the  = isqrt(t2*alt**_sage_const_2 )
    oot  = the * alt
    the *= t2
    tos  = the * s21
    X = _sage_const_2  * (tos-the) * oot
    Y = (_sage_const_1 -s2) * oot

    if not qpositive(tos): Y = -Y
    assert qpositive(X*Y)

    return X,Y

def test():
    P = _sage_const_2 *M.random_point()
    X,Y = M_to_E(P)
    s = decaf_encode_from_E(X,Y)
    assert s == decaf_encode_from_E_c(X,Y)
    XX,YY = decaf_decode_to_E(s)
    XX2,YY2 = decaf_decode_to_E_c(s)
    assert is_rotation((X,Y),(XX,YY))
    assert is_rotation((X,Y),(XX2,YY2))

P = _sage_const_2 *M.random_point()
X,Y = M_to_E(P)
print("the X", X)
print("the Y", Y)
print decaf_encode_from_E(X, Y)
print M.base_ring()
print M.weierstrass_p(prec=_sage_const_10 )
x = _sage_const_117812161263436946737282484343310064665180535357016373416879082147939404277809514858788439644911793978499419995990477371552926308078495 
y = _sage_const_19 
X1, Y1 = M_to_X(x, y)
print("the x", X1)
print("the y", Y1)

