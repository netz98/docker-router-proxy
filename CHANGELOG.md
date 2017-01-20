# Changelog

**2.0.1**

*2017-01-20*

* Do not use the docker process pre-filter as this causes trouble in some edge-cases 
  * As the results are cached anyways, this doesn't have a performance impact
* Added some additional debuggin output for container resolving

**2.0.0**

*2016-12-22*

* BC break:
  * Only route web traffic (*port 80 and 443*)

**1.0.0**

* First release
