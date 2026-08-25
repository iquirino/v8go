#ifndef V8GO_SYMBOL_H
#define V8GO_SYMBOL_H

#include "errors.h"

#ifdef __cplusplus
namespace v8 {
class Isolate;
}
typedef v8::Isolate v8Isolate;
extern "C" {
#else
typedef struct v8Isolate v8Isolate;
#endif
typedef v8Isolate* IsolatePtr;

typedef struct m_value m_value;
typedef m_value* ValuePtr;

typedef enum {
  SYMBOL_ASYNC_ITERATOR = 1,
  SYMBOL_HAS_INSTANCE,
  SYMBOL_IS_CONCAT_SPREADABLE,
  SYMBOL_ITERATOR,
  SYMBOL_MATCH,
  SYMBOL_REPLACE,
  SYMBOL_SEARCH,
  SYMBOL_SPLIT,
  SYMBOL_TO_PRIMITIVE,
  SYMBOL_TO_STRING_TAG,
  SYMBOL_UNSCOPABLES,
} SymbolIndex;

ValuePtr BuiltinSymbol(IsolatePtr iso_ptr, SymbolIndex idx);
const char* SymbolDescription(ValuePtr ptr);

extern int ObjectSetPrivate(ValuePtr ptr, const char* key, int key_len, ValuePtr val_ptr);
extern RtnValue ObjectGetPrivate(ValuePtr ptr, const char* key, int key_len);
extern int ObjectHasPrivate(ValuePtr ptr, const char* key, int key_len);
extern int ObjectDeletePrivate(ValuePtr ptr, const char* key, int key_len);

#ifdef __cplusplus
}
#endif
#endif
