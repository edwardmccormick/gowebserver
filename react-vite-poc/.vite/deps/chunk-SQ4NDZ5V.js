import {
  canUseDOM_default
} from "./chunk-S2TLU4L2.js";
import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/dom-helpers/esm/contains.js
function contains(context, node) {
  if (context.contains) return context.contains(node);
  if (context.compareDocumentPosition) return context === node || !!(context.compareDocumentPosition(node) & 16);
}

// node_modules/@restart/ui/esm/useWindow.js
var import_react = __toESM(require_react());
var Context = (0, import_react.createContext)(canUseDOM_default ? window : void 0);
var WindowProvider = Context.Provider;
function useWindow() {
  return (0, import_react.useContext)(Context);
}

export {
  contains,
  useWindow
};
//# sourceMappingURL=chunk-SQ4NDZ5V.js.map
