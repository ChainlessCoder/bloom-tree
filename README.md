[![](https://img.shields.io/badge/made%20by-Bloom%20Lab-blue.svg?style=flat-square)](https://bloomlab.io)
[![Build Status](https://travis-ci.com/labbloom/bloom-tree.svg?token=KzkBQ6duyh2GgqS9Be5J&branch=master)](https://travis-ci.com/labbloom/bloom-tree)
[![codecov](https://codecov.io/gh/labbloom/bloom-tree/branch/master/graph/badge.svg?token=xLnQTvQe2W)](https://codecov.io/gh/labbloom/bloom-tree)

# The Bloom Tree
The Bloom Tree combines the idea of bloom filters with that of merkle trees. 
In the standard bloom filter, we are interested to verify the presence, or absence of element(s) in a set. 
In the case of the  Bloom Tree, we are interested to check and transmit the presence, or absence of an element in a secure and bandwidth efficient way to another party. 
Instead of sending the whole bloom filter to a receiver, we only send a small multiproof.

<img src="https://github.com/labbloom/bloom-tree/blob/master/img/bloom-tree.png" class="center" width="900" height="380">
