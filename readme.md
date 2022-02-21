
<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->


# CMS [![Makefile CI](https://github.com/crikke/ffcms/actions/workflows/makefile.yml/badge.svg?branch=master)](https://github.com/crikke/ffcms/actions/workflows/makefile.yml)
V1 board can be found [here](https://github.com/crikke/ffcms/projects/1)
## Overview
A headless CMS. 
use chi instead of gin

### Content Delivery Api
Decoupled from ContentManagementApi
Responsible for fetching content
Read optimized

### Content Management Api
Responsible for creating content & ContentTypes
