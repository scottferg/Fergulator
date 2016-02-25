#include <stdlib.h>
#include <GL/glew.h>

#undef GLEW_GET_FUN
#define GLEW_GET_FUN(x) (*x)