"use client";
import {
  useWillUnmount
} from "./chunk-OMFDF2U2.js";
import {
  useCommittedRef_default,
  useEventCallback
} from "./chunk-PJQH5SKO.js";
import {
  forEach,
  map
} from "./chunk-2A7GM3I4.js";
import {
  TransitionWrapper_default,
  transitionEndListener,
  triggerBrowserReflow
} from "./chunk-T2ES4POL.js";
import "./chunk-S2TLU4L2.js";
import "./chunk-QOVZGY7A.js";
import "./chunk-L3ZLE5PF.js";
import {
  Anchor_default
} from "./chunk-EKOP4GHW.js";
import "./chunk-ZEXTMIKL.js";
import "./chunk-MKS2UZZH.js";
import {
  useUncontrolled
} from "./chunk-BGOLUZHQ.js";
import "./chunk-EEJ6RG5R.js";
import {
  require_classnames,
  require_jsx_runtime,
  useBootstrapPrefix,
  useIsRTL
} from "./chunk-5TWWBDIN.js";
import {
  __toESM,
  require_react
} from "./chunk-5WQJO2FO.js";

// node_modules/@restart/hooks/esm/useUpdateEffect.js
var import_react = __toESM(require_react());
function useUpdateEffect(fn, deps) {
  const isFirst = (0, import_react.useRef)(true);
  (0, import_react.useEffect)(() => {
    if (isFirst.current) {
      isFirst.current = false;
      return;
    }
    return fn();
  }, deps);
}
var useUpdateEffect_default = useUpdateEffect;

// node_modules/@restart/hooks/esm/useTimeout.js
var import_react3 = __toESM(require_react());

// node_modules/@restart/hooks/esm/useMounted.js
var import_react2 = __toESM(require_react());
function useMounted() {
  const mounted = (0, import_react2.useRef)(true);
  const isMounted = (0, import_react2.useRef)(() => mounted.current);
  (0, import_react2.useEffect)(() => {
    mounted.current = true;
    return () => {
      mounted.current = false;
    };
  }, []);
  return isMounted.current;
}

// node_modules/@restart/hooks/esm/useTimeout.js
var MAX_DELAY_MS = 2 ** 31 - 1;
function setChainedTimeout(handleRef, fn, timeoutAtMs) {
  const delayMs = timeoutAtMs - Date.now();
  handleRef.current = delayMs <= MAX_DELAY_MS ? setTimeout(fn, delayMs) : setTimeout(() => setChainedTimeout(handleRef, fn, timeoutAtMs), MAX_DELAY_MS);
}
function useTimeout() {
  const isMounted = useMounted();
  const handleRef = (0, import_react3.useRef)();
  useWillUnmount(() => clearTimeout(handleRef.current));
  return (0, import_react3.useMemo)(() => {
    const clear = () => clearTimeout(handleRef.current);
    function set(fn, delayMs = 0) {
      if (!isMounted()) return;
      clear();
      if (delayMs <= MAX_DELAY_MS) {
        handleRef.current = setTimeout(fn, delayMs);
      } else {
        setChainedTimeout(handleRef, fn, Date.now() + delayMs);
      }
    }
    return {
      set,
      clear,
      handleRef
    };
  }, []);
}

// node_modules/react-bootstrap/esm/Carousel.js
var import_classnames3 = __toESM(require_classnames());
var React3 = __toESM(require_react());
var import_react4 = __toESM(require_react());

// node_modules/react-bootstrap/esm/CarouselCaption.js
var React = __toESM(require_react());
var import_classnames = __toESM(require_classnames());
var import_jsx_runtime = __toESM(require_jsx_runtime());
var CarouselCaption = React.forwardRef(({
  className,
  bsPrefix,
  as: Component = "div",
  ...props
}, ref) => {
  bsPrefix = useBootstrapPrefix(bsPrefix, "carousel-caption");
  return (0, import_jsx_runtime.jsx)(Component, {
    ref,
    className: (0, import_classnames.default)(className, bsPrefix),
    ...props
  });
});
CarouselCaption.displayName = "CarouselCaption";
var CarouselCaption_default = CarouselCaption;

// node_modules/react-bootstrap/esm/CarouselItem.js
var import_classnames2 = __toESM(require_classnames());
var React2 = __toESM(require_react());
var import_jsx_runtime2 = __toESM(require_jsx_runtime());
var CarouselItem = React2.forwardRef(({
  // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
  as: Component = "div",
  bsPrefix,
  className,
  ...props
}, ref) => {
  const finalClassName = (0, import_classnames2.default)(className, useBootstrapPrefix(bsPrefix, "carousel-item"));
  return (0, import_jsx_runtime2.jsx)(Component, {
    ref,
    ...props,
    className: finalClassName
  });
});
CarouselItem.displayName = "CarouselItem";
var CarouselItem_default = CarouselItem;

// node_modules/react-bootstrap/esm/Carousel.js
var import_jsx_runtime3 = __toESM(require_jsx_runtime());
var import_jsx_runtime4 = __toESM(require_jsx_runtime());
var import_jsx_runtime5 = __toESM(require_jsx_runtime());
var SWIPE_THRESHOLD = 40;
function isVisible(element) {
  if (!element || !element.style || !element.parentNode || !element.parentNode.style) {
    return false;
  }
  const elementStyle = getComputedStyle(element);
  return elementStyle.display !== "none" && elementStyle.visibility !== "hidden" && getComputedStyle(element.parentNode).display !== "none";
}
var Carousel = (
  // eslint-disable-next-line react/display-name
  React3.forwardRef(({
    defaultActiveIndex = 0,
    ...uncontrolledProps
  }, ref) => {
    const {
      // Need to define the default "as" during prop destructuring to be compatible with styled-components github.com/react-bootstrap/react-bootstrap/issues/3595
      as: Component = "div",
      bsPrefix,
      slide = true,
      fade = false,
      controls = true,
      indicators = true,
      indicatorLabels = [],
      activeIndex,
      onSelect,
      onSlide,
      onSlid,
      interval = 5e3,
      keyboard = true,
      onKeyDown,
      pause = "hover",
      onMouseOver,
      onMouseOut,
      wrap = true,
      touch = true,
      onTouchStart,
      onTouchMove,
      onTouchEnd,
      prevIcon = (0, import_jsx_runtime3.jsx)("span", {
        "aria-hidden": "true",
        className: "carousel-control-prev-icon"
      }),
      prevLabel = "Previous",
      nextIcon = (0, import_jsx_runtime3.jsx)("span", {
        "aria-hidden": "true",
        className: "carousel-control-next-icon"
      }),
      nextLabel = "Next",
      variant,
      className,
      children,
      ...props
    } = useUncontrolled({
      defaultActiveIndex,
      ...uncontrolledProps
    }, {
      activeIndex: "onSelect"
    });
    const prefix = useBootstrapPrefix(bsPrefix, "carousel");
    const isRTL = useIsRTL();
    const nextDirectionRef = (0, import_react4.useRef)(null);
    const [direction, setDirection] = (0, import_react4.useState)("next");
    const [paused, setPaused] = (0, import_react4.useState)(false);
    const [isSliding, setIsSliding] = (0, import_react4.useState)(false);
    const [renderedActiveIndex, setRenderedActiveIndex] = (0, import_react4.useState)(activeIndex || 0);
    (0, import_react4.useEffect)(() => {
      if (!isSliding && activeIndex !== renderedActiveIndex) {
        if (nextDirectionRef.current) {
          setDirection(nextDirectionRef.current);
        } else {
          setDirection((activeIndex || 0) > renderedActiveIndex ? "next" : "prev");
        }
        if (slide) {
          setIsSliding(true);
        }
        setRenderedActiveIndex(activeIndex || 0);
      }
    }, [activeIndex, isSliding, renderedActiveIndex, slide]);
    (0, import_react4.useEffect)(() => {
      if (nextDirectionRef.current) {
        nextDirectionRef.current = null;
      }
    });
    let numChildren = 0;
    let activeChildInterval;
    forEach(children, (child, index) => {
      ++numChildren;
      if (index === activeIndex) {
        activeChildInterval = child.props.interval;
      }
    });
    const activeChildIntervalRef = useCommittedRef_default(activeChildInterval);
    const prev = (0, import_react4.useCallback)((event) => {
      if (isSliding) {
        return;
      }
      let nextActiveIndex = renderedActiveIndex - 1;
      if (nextActiveIndex < 0) {
        if (!wrap) {
          return;
        }
        nextActiveIndex = numChildren - 1;
      }
      nextDirectionRef.current = "prev";
      onSelect == null || onSelect(nextActiveIndex, event);
    }, [isSliding, renderedActiveIndex, onSelect, wrap, numChildren]);
    const next = useEventCallback((event) => {
      if (isSliding) {
        return;
      }
      let nextActiveIndex = renderedActiveIndex + 1;
      if (nextActiveIndex >= numChildren) {
        if (!wrap) {
          return;
        }
        nextActiveIndex = 0;
      }
      nextDirectionRef.current = "next";
      onSelect == null || onSelect(nextActiveIndex, event);
    });
    const elementRef = (0, import_react4.useRef)();
    (0, import_react4.useImperativeHandle)(ref, () => ({
      element: elementRef.current,
      prev,
      next
    }));
    const nextWhenVisible = useEventCallback(() => {
      if (!document.hidden && isVisible(elementRef.current)) {
        if (isRTL) {
          prev();
        } else {
          next();
        }
      }
    });
    const slideDirection = direction === "next" ? "start" : "end";
    useUpdateEffect_default(() => {
      if (slide) {
        return;
      }
      onSlide == null || onSlide(renderedActiveIndex, slideDirection);
      onSlid == null || onSlid(renderedActiveIndex, slideDirection);
    }, [renderedActiveIndex]);
    const orderClassName = `${prefix}-item-${direction}`;
    const directionalClassName = `${prefix}-item-${slideDirection}`;
    const handleEnter = (0, import_react4.useCallback)((node) => {
      triggerBrowserReflow(node);
      onSlide == null || onSlide(renderedActiveIndex, slideDirection);
    }, [onSlide, renderedActiveIndex, slideDirection]);
    const handleEntered = (0, import_react4.useCallback)(() => {
      setIsSliding(false);
      onSlid == null || onSlid(renderedActiveIndex, slideDirection);
    }, [onSlid, renderedActiveIndex, slideDirection]);
    const handleKeyDown = (0, import_react4.useCallback)((event) => {
      if (keyboard && !/input|textarea/i.test(event.target.tagName)) {
        switch (event.key) {
          case "ArrowLeft":
            event.preventDefault();
            if (isRTL) {
              next(event);
            } else {
              prev(event);
            }
            return;
          case "ArrowRight":
            event.preventDefault();
            if (isRTL) {
              prev(event);
            } else {
              next(event);
            }
            return;
          default:
        }
      }
      onKeyDown == null || onKeyDown(event);
    }, [keyboard, onKeyDown, prev, next, isRTL]);
    const handleMouseOver = (0, import_react4.useCallback)((event) => {
      if (pause === "hover") {
        setPaused(true);
      }
      onMouseOver == null || onMouseOver(event);
    }, [pause, onMouseOver]);
    const handleMouseOut = (0, import_react4.useCallback)((event) => {
      setPaused(false);
      onMouseOut == null || onMouseOut(event);
    }, [onMouseOut]);
    const touchStartXRef = (0, import_react4.useRef)(0);
    const touchDeltaXRef = (0, import_react4.useRef)(0);
    const touchUnpauseTimeout = useTimeout();
    const handleTouchStart = (0, import_react4.useCallback)((event) => {
      touchStartXRef.current = event.touches[0].clientX;
      touchDeltaXRef.current = 0;
      if (pause === "hover") {
        setPaused(true);
      }
      onTouchStart == null || onTouchStart(event);
    }, [pause, onTouchStart]);
    const handleTouchMove = (0, import_react4.useCallback)((event) => {
      if (event.touches && event.touches.length > 1) {
        touchDeltaXRef.current = 0;
      } else {
        touchDeltaXRef.current = event.touches[0].clientX - touchStartXRef.current;
      }
      onTouchMove == null || onTouchMove(event);
    }, [onTouchMove]);
    const handleTouchEnd = (0, import_react4.useCallback)((event) => {
      if (touch) {
        const touchDeltaX = touchDeltaXRef.current;
        if (Math.abs(touchDeltaX) > SWIPE_THRESHOLD) {
          if (touchDeltaX > 0) {
            prev(event);
          } else {
            next(event);
          }
        }
      }
      if (pause === "hover") {
        touchUnpauseTimeout.set(() => {
          setPaused(false);
        }, interval || void 0);
      }
      onTouchEnd == null || onTouchEnd(event);
    }, [touch, pause, prev, next, touchUnpauseTimeout, interval, onTouchEnd]);
    const shouldPlay = interval != null && !paused && !isSliding;
    const intervalHandleRef = (0, import_react4.useRef)();
    (0, import_react4.useEffect)(() => {
      var _ref, _activeChildIntervalR;
      if (!shouldPlay) {
        return void 0;
      }
      const nextFunc = isRTL ? prev : next;
      intervalHandleRef.current = window.setInterval(document.visibilityState ? nextWhenVisible : nextFunc, (_ref = (_activeChildIntervalR = activeChildIntervalRef.current) != null ? _activeChildIntervalR : interval) != null ? _ref : void 0);
      return () => {
        if (intervalHandleRef.current !== null) {
          clearInterval(intervalHandleRef.current);
        }
      };
    }, [shouldPlay, prev, next, activeChildIntervalRef, interval, nextWhenVisible, isRTL]);
    const indicatorOnClicks = (0, import_react4.useMemo)(() => indicators && Array.from({
      length: numChildren
    }, (_, index) => (event) => {
      onSelect == null || onSelect(index, event);
    }), [indicators, numChildren, onSelect]);
    return (0, import_jsx_runtime4.jsxs)(Component, {
      ref: elementRef,
      ...props,
      onKeyDown: handleKeyDown,
      onMouseOver: handleMouseOver,
      onMouseOut: handleMouseOut,
      onTouchStart: handleTouchStart,
      onTouchMove: handleTouchMove,
      onTouchEnd: handleTouchEnd,
      className: (0, import_classnames3.default)(className, prefix, slide && "slide", fade && `${prefix}-fade`, variant && `${prefix}-${variant}`),
      children: [indicators && (0, import_jsx_runtime3.jsx)("div", {
        className: `${prefix}-indicators`,
        children: map(children, (_, index) => (0, import_jsx_runtime3.jsx)("button", {
          type: "button",
          "data-bs-target": "",
          "aria-label": indicatorLabels != null && indicatorLabels.length ? indicatorLabels[index] : `Slide ${index + 1}`,
          className: index === renderedActiveIndex ? "active" : void 0,
          onClick: indicatorOnClicks ? indicatorOnClicks[index] : void 0,
          "aria-current": index === renderedActiveIndex
        }, index))
      }), (0, import_jsx_runtime3.jsx)("div", {
        className: `${prefix}-inner`,
        children: map(children, (child, index) => {
          const isActive = index === renderedActiveIndex;
          return slide ? (0, import_jsx_runtime3.jsx)(TransitionWrapper_default, {
            in: isActive,
            onEnter: isActive ? handleEnter : void 0,
            onEntered: isActive ? handleEntered : void 0,
            addEndListener: transitionEndListener,
            children: (status, innerProps) => React3.cloneElement(child, {
              ...innerProps,
              className: (0, import_classnames3.default)(child.props.className, isActive && status !== "entered" && orderClassName, (status === "entered" || status === "exiting") && "active", (status === "entering" || status === "exiting") && directionalClassName)
            })
          }) : React3.cloneElement(child, {
            className: (0, import_classnames3.default)(child.props.className, isActive && "active")
          });
        })
      }), controls && (0, import_jsx_runtime4.jsxs)(import_jsx_runtime5.Fragment, {
        children: [(wrap || activeIndex !== 0) && (0, import_jsx_runtime4.jsxs)(Anchor_default, {
          className: `${prefix}-control-prev`,
          onClick: prev,
          children: [prevIcon, prevLabel && (0, import_jsx_runtime3.jsx)("span", {
            className: "visually-hidden",
            children: prevLabel
          })]
        }), (wrap || activeIndex !== numChildren - 1) && (0, import_jsx_runtime4.jsxs)(Anchor_default, {
          className: `${prefix}-control-next`,
          onClick: next,
          children: [nextIcon, nextLabel && (0, import_jsx_runtime3.jsx)("span", {
            className: "visually-hidden",
            children: nextLabel
          })]
        })]
      })]
    });
  })
);
Carousel.displayName = "Carousel";
var Carousel_default = Object.assign(Carousel, {
  Caption: CarouselCaption_default,
  Item: CarouselItem_default
});
export {
  Carousel_default as default
};
//# sourceMappingURL=react-bootstrap_Carousel.js.map
