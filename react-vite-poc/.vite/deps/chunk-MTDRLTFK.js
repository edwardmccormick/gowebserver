import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/@restart/ui/esm/utils.js
var React = __toESM(require_react());
function isEscKey(e) {
  return e.code === "Escape" || e.keyCode === 27;
}
function getReactVersion() {
  const parts = React.version.split(".");
  return {
    major: +parts[0],
    minor: +parts[1],
    patch: +parts[2]
  };
}
function getChildRef(element) {
  if (!element || typeof element === "function") {
    return null;
  }
  const {
    major
  } = getReactVersion();
  const childRef = major >= 19 ? element.props.ref : element.ref;
  return childRef;
}

export {
  isEscKey,
  getChildRef
};
//# sourceMappingURL=chunk-MTDRLTFK.js.map
