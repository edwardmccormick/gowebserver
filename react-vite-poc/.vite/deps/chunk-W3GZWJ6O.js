// node_modules/dom-helpers/esm/querySelectorAll.js
var toArray = Function.prototype.bind.call(Function.prototype.call, [].slice);
function qsa(element, selector) {
  return toArray(element.querySelectorAll(selector));
}

export {
  qsa
};
//# sourceMappingURL=chunk-W3GZWJ6O.js.map
