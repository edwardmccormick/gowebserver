import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/@restart/ui/esm/SelectableContext.js
var React = __toESM(require_react());
var SelectableContext = React.createContext(null);
var makeEventKey = (eventKey, href = null) => {
  if (eventKey != null) return String(eventKey);
  return href || null;
};
var SelectableContext_default = SelectableContext;

export {
  makeEventKey,
  SelectableContext_default
};
//# sourceMappingURL=chunk-25OZPJTB.js.map
