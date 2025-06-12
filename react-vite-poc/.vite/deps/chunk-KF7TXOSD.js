import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/@restart/ui/node_modules/@restart/hooks/esm/useForceUpdate.js
var import_react = __toESM(require_react());
function useForceUpdate() {
  const [, dispatch] = (0, import_react.useReducer)((revision) => revision + 1, 0);
  return dispatch;
}

export {
  useForceUpdate
};
//# sourceMappingURL=chunk-KF7TXOSD.js.map
