CMAKE_MINIMUM_REQUIRED (VERSION 2.8)
PROJECT (query)

# This project is intended to be built as part of a Couchbase Server build,
# but for nominal support building stand-alone, the following variables
# need to be specified:
#
# CMAKE_MODULE_PATH - should contain a directory containing
#      FindCouchbaseGo.cmake, eg. tlm/cmake/Modules
#
# CMAKE_INSTALL_PREFIX - should point to a directory where the bin/ output
#      directory will be placed
#
# GODEPSDIR - should point to a Go workspace directory containing all
#      transitive Go dependencies

INCLUDE (FindCouchbaseGo)

GoYacc (TARGET n1ql-yacc YFILE "${PROJECT_SOURCE_DIR}/parser/n1ql/n1ql.y")

GoInstall (TARGET cbq-engine PACKAGE github.com/couchbaselabs/query/server/cbq-engine
  GOPATH "${PROJECT_SOURCE_DIR}/../../../.." "${GODEPSDIR}"
  DEPENDS n1ql-yacc INSTALL_PATH bin OUTPUT cbq-engine)

GoInstall (TARGET cbq PACKAGE github.com/couchbaselabs/query/shell/cbq
  GOPATH "${PROJECT_SOURCE_DIR}/../../../.." "${GODEPSDIR}"
  DEPENDS n1ql-yacc cbq-engine INSTALL_PATH bin OUTPUT cbq)
