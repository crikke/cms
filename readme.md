
<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->


# CMS
## Overview
A headless CMS. 
Reliability, scalability, speed.

### Content Delivery Api
Decoupled from ContentManagementApi
Responsible for fetching content
Read optimized

### Content Management Api
Responsible for creating content & ContentTypes

## Choosing a Database
Cassandra vs MongoDB
Cassandra has multiple master nodes so write speeds are improved
Mongodb has a single primary node that is used for writes and multiple secondary nodes which are used for read.

> A good guess is that the CMS see more reads than write, so mongodb can work

## Authorization

