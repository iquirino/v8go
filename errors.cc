#include <stdlib.h>
#include <sstream>

#include "deps/include/v8-context.h"
#include "deps/include/v8-exception.h"
#include "deps/include/v8-external.h"
#include "deps/include/v8-message.h"
#include "deps/include/v8-primitive.h"

#include "context.h"
#include "errors.h"
#include "utils.h"
#include "value.h"

using namespace v8;

RtnError ExceptionError(TryCatch& try_catch, Isolate* iso, Local<Context> ctx) {
  HandleScope handle_scope(iso);

  RtnError rtn = {nullptr, nullptr, nullptr, nullptr};

  if (try_catch.HasTerminated()) {
    rtn.msg =
        CopyString("ExecutionTerminated: script execution has been terminated");
    return rtn;
  }

  Local<Value> exception = try_catch.Exception();

  String::Utf8Value exceptionStr(iso, exception);
  rtn.msg = CopyString(exceptionStr);

  Local<Message> msg = try_catch.Message();
  if (!msg.IsEmpty()) {
    String::Utf8Value origin(iso, msg->GetScriptOrigin().ResourceName());
    std::ostringstream sb;
    sb << *origin;
    Maybe<int> line = try_catch.Message()->GetLineNumber(ctx);
    if (line.IsJust()) {
      sb << ":" << line.ToChecked();
    }
    Maybe<int> start = try_catch.Message()->GetStartColumn(ctx);
    if (start.IsJust()) {
      sb << ":"
         << start.ToChecked() + 1;  // + 1 to match output from stack trace
    }
    rtn.location = CopyString(sb.str());
  }

  Local<Value> mstack;
  if (try_catch.StackTrace(ctx).ToLocal(&mstack)) {
    String::Utf8Value stack(iso, mstack);
    rtn.stack = CopyString(stack);
  }

  // Track the exception value in the context so it's freed on context close.
  if (!exception.IsEmpty() && !ctx.IsEmpty()) {
    Local<Data> embedder_data = ctx->GetEmbedderDataV2(2);
    if (!embedder_data.IsEmpty() && embedder_data->IsValue()) {
      Local<Value> embedder = Local<Value>::Cast(embedder_data);
      if (embedder->IsExternal()) {
        m_ctx* m_context = (m_ctx*)embedder.As<External>()->Value(
            kExternalPointerTypeTagDefault);
        if (m_context != nullptr) {
          rtn.exception_value = track_value(m_context, exception);
        }
      }
    }
  }

  return rtn;
}

void ErrorRelease(RtnError err) {
  free(err.msg);
  free(err.location);
  free(err.stack);
}
