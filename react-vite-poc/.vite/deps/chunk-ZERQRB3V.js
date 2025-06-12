import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/@restart/hooks/esm/useIsomorphicEffect.js
var import_react = __toESM(require_react());
var isReactNative = typeof global !== "undefined" && // @ts-ignore
global.navigator && // @ts-ignore
global.navigator.product === "ReactNative";
var isDOM = typeof document !== "undefined";
var useIsomorphicEffect_default = isDOM || isReactNative ? import_react.useLayoutEffect : import_react.useEffect;

export {
  useIsomorphicEffect_default
};
//# sourceMappingURL=chunk-ZERQRB3V.js.map
