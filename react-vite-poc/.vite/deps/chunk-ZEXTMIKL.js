import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/@restart/ui/node_modules/@restart/hooks/esm/useEventCallback.js
var import_react2 = __toESM(require_react());

// node_modules/@restart/ui/node_modules/@restart/hooks/esm/useCommittedRef.js
var import_react = __toESM(require_react());
function useCommittedRef(value) {
  const ref = (0, import_react.useRef)(value);
  (0, import_react.useEffect)(() => {
    ref.current = value;
  }, [value]);
  return ref;
}
var useCommittedRef_default = useCommittedRef;

// node_modules/@restart/ui/node_modules/@restart/hooks/esm/useEventCallback.js
function useEventCallback(fn) {
  const ref = useCommittedRef_default(fn);
  return (0, import_react2.useCallback)(function(...args) {
    return ref.current && ref.current(...args);
  }, [ref]);
}

// node_modules/@restart/ui/node_modules/@restart/hooks/esm/useMounted.js
var import_react3 = __toESM(require_react());
function useMounted() {
  const mounted = (0, import_react3.useRef)(true);
  const isMounted = (0, import_react3.useRef)(() => mounted.current);
  (0, import_react3.useEffect)(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
    };
  }, []);
  return isMounted.current;
}

// node_modules/@restart/ui/node_modules/@restart/hooks/esm/usePrevious.js
var import_react4 = __toESM(require_react());
function usePrevious(value) {
  const ref = (0, import_react4.useRef)(null);
  (0, import_react4.useEffect)(() => {
    ref.current = value;
  });
  return ref.current;
}

// node_modules/@restart/ui/node_modules/@restart/hooks/esm/useIsomorphicEffect.js
var import_react5 = __toESM(require_react());
var isReactNative = typeof global !== "undefined" && // @ts-ignore
global.navigator && // @ts-ignore
global.navigator.product === "ReactNative";
var isDOM = typeof document !== "undefined";
var useIsomorphicEffect_default = isDOM || isReactNative ? import_react5.useLayoutEffect : import_react5.useEffect;

export {
  useEventCallback,
  useMounted,
  usePrevious,
  useIsomorphicEffect_default
};
//# sourceMappingURL=chunk-ZEXTMIKL.js.map
