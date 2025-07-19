(function (win) {
	var amplayer = {
		'options': {
			// player options
			"autoPlay": true,
			"screenshot": false,
			"airplay": true,
			"chromecast": true,
			"live": false, // 直播模式
			"logo": "",

			// video options
			"urls": "",
			"pics": "",

			// other options
			"autoSkip": true, //自动跳过片头和片尾
			"autoSkipAd": true, //自动跳过广告
			"autoNext": true,
			"debug": false,
			"trys": 0, // 试看时长 seconds
			"seek": 0, // 跳过时长 seconds
			"take": "",
			"seq": "",
			"jump": "",
			"p2pAppId": "",
			"p2pEngine": "",
			"p2pConfig": {},
			"srt2vttConvertApi": "",
			"container": "",
			"defaultType": "customHls",
			"defaultExtName": ".m3u8",
			"touchVideoChangeProgress": false,
			"listeners": {}
		},
		'secure': win.location.protocol == 'https:',
		'elemPrefix': function (notPrefix) {
			return (amplayer.options.container ? amplayer.options.container + (!notPrefix?' ':'') : '');
		},
		'player': {
			'torrentAdd': function (torrent, video, player) {
				player.notice(App.i18n.CONNECTED, 5000);
				//console.dir(torrent.files);
				var reVideo = new RegExp('\\.(mp4)$', 'i');
				var file = torrent.files.find(function (file) {
					console.log('[torrent] adding file: ', file.path || file.name);
					return reVideo.test(file.name);
				});
				var renderOptions = {
					autoplay: player.options.autoplay
				};
				if (torrent.urlList && torrent.urlList.length > 0) {
					var reCover = new RegExp('poster\\.(jp[e]?g|png|gif|webp|svg)$', 'i');
					var cover = torrent.files.find(function (file) {
						return reCover.test(file.name);
					});
					if (cover) {
						cover = torrent.urlList[0] + cover.path;
						console.log('[torrent] adding cover: ', cover);
						$(amplayer.elemPrefix() + 'video').attr('poster', cover);
					}
					var reSubtitleVTT = new RegExp('\\.(vtt)$', 'i');
					var subtitle = torrent.files.find(function (file) {
						return reSubtitleVTT.test(file.name);
					});
					var subtitleType = 'webvtt', subtitleURL = '';
					if (!subtitle) {
						if (amplayer.options.srt2vttConvertApi) {
							var reSubtitleSRT = new RegExp('\\.(srt)$', 'i');
							subtitle = torrent.files.find(function (file) {
								return reSubtitleSRT.test(file.name);
							});
							if (subtitle) {
								subtitleURL = torrent.urlList[0] + subtitle.path;
								subtitleURL = amplayer.options.srt2vttConvertApi + encodeURIComponent(subtitleURL);
							}
						}
					} else {
						subtitleURL = torrent.urlList[0] + subtitle.path;
					}
					if (subtitleURL) {
						console.log('[torrent] adding subtitle: ', subtitleURL);
						player.options.subtitle = {
							type: subtitleType,
							url: subtitleURL,
							fontSize: '25px',
							bottom: '10%',
							color: '#b7daff'
						}
						player.initSubtitle(player.options.subtitle);
					}
				}
				if (!file) {
					player.notice(App.i18n.PLAY_FAILED + ': ' + App.i18n.NOT_FOUND_MEDIA, 5000);
					return;
				}
				player.notice(App.i18n.READY_PLAY, 50000);
				file.renderTo(video, renderOptions, function () {
					player.container.classList.remove('dplayer-loading');
				});
			},
			'customType': {
				'customWebTorrent': function (video, player) {
					//测试种子: https://webtorrent.io/torrents/sintel.torrent
					player.container.classList.add('dplayer-loading');
					var client = new WebTorrent();
					var torrentId = video.src;
					video.torrentId = torrentId;
					video.src = '';
					video.preload = 'metadata';
					video.addEventListener('durationchange', function () {
						player.container.classList.remove('dplayer-loading');
					}, { once: true });
					//console.log("\n %c Added torrentId %c "+torrentId+" \n\n","color: #fadfa3; background: #030307; padding:5px 0;","background: #fadfa3; padding:5px 0;");
					player.notice(App.i18n.CONNECTING, 5000);
					var _noticeFixerWatcher = win.setInterval(function () {
						if ($(amplayer.elemPrefix() + '.dplayer-notice').text() == '视频加载失败') {
							player.notice(App.i18n.CONNECTING, 5000);
							win.clearInterval(_noticeFixerWatcher);
						}
					}, 50);
					var _torrentWatcher = null;
					var clearNFW = function () {
						if (_noticeFixerWatcher) {
							clearInterval(_noticeFixerWatcher);
							_noticeFixerWatcher = null;
						}
					};
					var clearTW = function () {
						if (_torrentWatcher) {
							clearInterval(_torrentWatcher);
							_torrentWatcher = null;
						}
					};
					client.add(torrentId, function (torrent) {
						console.log('[torrent] Client is downloading:', torrent.infoHash);
						amplayer.player.torrentAdd(torrent, video, player);

						var onProgress = function () {
							clearNFW();
							var percent = Math.round(torrent.progress * 100 * 100) / 100;
							$(amplayer.elemPrefix() + '.line').html('速度 <span class="fa fa-long-arrow-down text-success"></span>' + App.formatBytes(torrent.downloadSpeed) + ' <span class="fa fa-long-arrow-up text-danger"></span>' + App.formatBytes(torrent.uploadSpeed) + ' 在线' + torrent.numPeers + 'NP');
							$(amplayer.elemPrefix() + '.peer').text('BT已开启');//torrent.timeRemaining
							var msg = '已加载' + percent + '%';
							if (torrent.downloaded) {
								msg += ' (' + App.formatBytes(torrent.downloaded) + ')';
							}
							if (torrent.length) {
								msg += ' 共' + App.formatBytes(torrent.length);
							}
							$(amplayer.elemPrefix() + '.load').text(msg);
						};

						_torrentWatcher = setInterval(onProgress, 5000);
						onProgress();
						// torrent.on('download', onProgress);
						// torrent.on('upload', onProgress);
						torrent.on('wire', function (wire, addr) {
							console.log('connected to peer with address ' + addr)
							//wire.use(MyExtension)
							clearNFW();
						});
						torrent.on('noPeers', function (announceType) {
							console.log('no peers');
						});
						torrent.on('error', function () {
							clearNFW();
						});
						torrent.on('done', function () {
							clearTW();
							clearNFW();
						})
					});
					$(amplayer.elemPrefix() + '#video').data('webTorrent', client);
					player.on('destroy', function () {
						clearTW();
						clearNFW();
						if (client) {
							if (client.get(torrentId)) client.remove(torrentId);
							client.destroy();
							client = null;
						}
						$(amplayer.elemPrefix() + '#video').data('webTorrent', null);
					});
					return client;
				},
				'customHls': function (video, player) {
					var vd = $(amplayer.elemPrefix() + '#video');
					if (vd.data('hls')) {
						var eventIndex = vd.data('eventIndex');
						player.trigger('destroy');
						if (eventIndex !== null) player.off('destroy', eventIndex);
					}
					var config = $.extend({
						debug: amplayer.options.debug,
						enableWorker: true,
						maxMaxBufferLength: 100, // seconds (default:600)
						maxBufferSize: 0, // bytes (default:60M)
						maxBufferLength: 15, // seconds (default:30)
						backBufferLength: 15,
						liveMaxBackBufferLength: 15,
						liveBackBufferLength: 15,
						liveSyncDurationCount: 1,
					}, player.options.pluginOptions.hls || {});
					var engine = null;
					if (amplayer.options.p2pEngine == 'p2p-media-loader') {
						if (p2pml.hlsjs.Engine.isSupported()) {
							engine = new p2pml.hlsjs.Engine();
							config.liveSyncDurationCount = 7 // To have at least 7 segments in queue
							config.loader = engine.createLoaderClass()
						} else {
							amplayer.options.p2pEngine = ''
						}
					} else {
						config.p2pConfig = {
							logLevel: 'error',
							live: amplayer.options.live,        // set to true in live mode
							//wsSignalerAddr: 'wss://opensignal.cdnbye.com',
							//announce: 'https://tracker.cdnbye.com/v1',
							appId: amplayer.options.p2pAppId || '', // 长度不操作30
							// Other p2pConfig options provided by CDNBye
							//tag: '', // 长度不超过20
							//appName: '', // 长度不超过30
							//token: 'free', // 长度不超过20
						}
						if (config.debug) config.p2pConfig.logLevel = 'debug';
						config.p2pConfig = $.extend(config.p2pConfig, amplayer.options.p2pConfig || {});
					}
					var hls = new Hls(config);
					if (engine) {
						p2pml.hlsjs.initHlsJsPlayer(hls);
					} else {
						engine = hls.p2pEngine || hls.engine; // hls.p2pEngine - cdnbye; hls.engine - raycdn
						if (!engine && P2PEngine && P2PEngine.isSupported()) {
							engine = new P2PEngine(hls, config.p2pConfig);
						}
					}
					hls.loadSource(video.src);
					hls.attachMedia(video);
					if (engine) {
						if (amplayer.options.p2pEngine == 'p2p-media-loader') {
							/*engine.on(p2pml.core.Events.SegmentLoaded, function(segment, peerId) {
								$('.load').text('加载0MB 共享0MB 加速0MB');
								$('.peer').text('P2P已开启');
								$('.line').text('在线1NP');
							});*/
							var peers = 0;
							engine.on(p2pml.core.Events.PeerConnect, function (peer) {
								peers++;
								$(amplayer.elemPrefix() + '.line').text('在线' + (peers + 1) + 'NP');
								$(amplayer.elemPrefix() + '.peer').text('P2P已开启');
							});
							engine.on(p2pml.core.Events.PeerClose, function (peerId) {
								peers--;
								if (peers <= 0) peers = 0;
								$(amplayer.elemPrefix() + '.line').text('在线' + (peers + 1) + 'NP');
							});
							var downloaded = { http: 0, p2p: 0 }, uploaded = { http: 0, p2p: 0 };
							engine.on(p2pml.core.Events.PieceBytesDownloaded, function (method, bytes, peerId) {
								downloaded[method] += bytes / 1024;
								$(amplayer.elemPrefix() + '.load').text('加载' + (downloaded.http / 1024).toFixed(2) + 'MB 共享' + (uploaded.p2p / 1024).toFixed(2) + 'MB 加速' + (downloaded.p2p / 1024).toFixed(2) + 'MB');
								downloaded.p2p > 0 ? $('.peer').text('P2P加速中') : $('.peer').text('P2P已开启');
							});
							engine.on(p2pml.core.Events.PieceBytesUploaded, function (method, bytes) {
								uploaded[method] += bytes / 1024;
								$(amplayer.elemPrefix() + '.load').text('加载' + (downloaded.http / 1024).toFixed(2) + 'MB 共享' + (uploaded.p2p / 1024).toFixed(2) + 'MB 加速' + (downloaded.p2p / 1024).toFixed(2) + 'MB');
								downloaded.p2p > 0 ? $('.peer').text('P2P加速中') : $('.peer').text('P2P已开启');
							});
						} else {
							engine.on('peerId', function (peerId) {
								$(amplayer.elemPrefix() + '.load').text('加载0MB 共享0MB 加速0MB');
								$(amplayer.elemPrefix() + '.peer').text('P2P已开启');
								$(amplayer.elemPrefix() + '.line').text('在线1NP');
							});
							engine.on('peers', function (peers) {
								$(amplayer.elemPrefix() + '.line').text('在线' + (peers.length + 1) + 'NP');
								$(amplayer.elemPrefix() + '.peer').text('P2P已开启');
							});
							engine.on('stats', function (data) {
								$(amplayer.elemPrefix() + '.load').text('加载' + (data.totalHTTPDownloaded / 1024).toFixed(2) + 'MB 共享' + (data.totalP2PUploaded / 1024).toFixed(2) + 'MB 加速' + (data.totalP2PDownloaded / 1024).toFixed(2) + 'MB');
								data.totalP2PDownloaded ? $(amplayer.elemPrefix() + '.peer').text('P2P加速中') : $(amplayer.elemPrefix() + '.peer').text('P2P已开启');
							});
						}
					}
					var recoverDecodingErrorDate, recoverSwapAudioCodecDate, recoverStartLoadDate;
					hls.on(Hls.Events.ERROR, function (event, data) {
						var msg = '';
						switch (data.type) {
							case 'mediaError':
								msg = '媒体错误';
								var autoRecover = amplayer.options.autoRecoverMediaError || false;
								if (autoRecover && data.fatal) {
									var now = (new Date()).getTime();
									if (!recoverDecodingErrorDate || now - recoverDecodingErrorDate > 3000) {
										recoverDecodingErrorDate = now;
										hls.recoverMediaError();
										return;
									}
									if (!recoverSwapAudioCodecDate || now - recoverSwapAudioCodecDate > 3000) {
										recoverSwapAudioCodecDate = now;
										hls.swapAudioCodec();
										hls.recoverMediaError();
										return;
									}
								}
								break;

							case 'networkError':
								msg = '网络错误';
								var autoRecover = amplayer.options.autoRecoverNetworkError || false;
								if (autoRecover && data.fatal) {
									var now = (new Date()).getTime();
									if (!recoverStartLoadDate || now - recoverStartLoadDate > 3000) {
										recoverStartLoadDate = now;
										hls.startLoad();
										return;
									}
								}
								break;

							default: msg = data.type;
						}
						msg += ': ';
						switch (data.details) {
							//mediaError
							case 'bufferFullError':
							case 'bufferSeekOverHole':
							case 'bufferNudgeOnStall':
							//networkError
							case 'fragLoadError'://离开页面时触发
								if (!data.fatal) return;
								msg += data.details;
								break;
							//networkError
							case 'manifestLoadError':
								msg += '媒体加载失败';
								break;
							case 'manifestLoadTimeOut':
								msg += '媒体加载超时';
								break;
							default:
								msg += data.details;
						}
						data.message = msg;
						console.error(msg);
						$(amplayer.elemPrefix() + '#video').trigger('error', data);
					});
					vd.data('hls', hls);
					var eventIndex = player.on('destroy', function () {
						if (hls) {
							hls.destroy();
							if (amplayer.options.debug) console.debug('hls destory.')
							hls = null;
						}
						if (engine) {
							if (typeof (engine.destroy) == 'function') {
								engine.destroy();
								if (amplayer.options.debug) console.debug('p2p engine destory.')
							}
							engine = null;
						}
						$(amplayer.elemPrefix() + '#video').data('hls', null);
					});
					vd.data('eventIndex', eventIndex);
					return hls;
				},
				'shakaDash': function (video, player) {
					var src = video.src;
					var playerShaka = new shaka.Player(video); // 将会修改 video.src
					playerShaka.load(src);
				}
			},
			'getType': function (urls) {
				var type = 'auto', urls = String(urls).split('#')[0].split('?')[0],
					pos = urls.lastIndexOf('.'), extName = pos > -1 ? urls.substring(pos).toLowerCase() : amplayer.options.defaultExtName;
				if (urls.substring(0, 7).toLowerCase() == 'magnet:' || extName == '.torrent') {
					type = 'customWebTorrent'; // webtorrent
				} else if (extName == '.m3u8') {
					type = 'customHls'; // hls
				} else if (extName == '.mpd') {
					type = 'dash'; // dash
				} else if (extName == '.flv') {
					type = 'flv'; // flv
				} else if (extName == '.mp4' || extName == '.mp3' || extName == '.webm' || extName == '.ogg' || extName == '.mkv') {
					type = 'normal';
				} else if (amplayer.options.defaultType) {
					type = amplayer.options.defaultType;
				}
				return type;
			},
			'eplayer': function (options) {
				var c = $.extend(amplayer.options, options || {});
				var type = amplayer.player.getType(c.urls);
				var ctn = amplayer.elemPrefix(true);
				var elem = ctn && $(ctn).length>0 ? $(ctn)[0] : null;
				if (!elem) elem = document.getElementById('video');
				var opts = {
					container: elem,
					autoplay: c.autoPlay,
					live: c.live,
					logo: c.logo,
					screenshot: c.screenshot,
					airplay: c.airplay,
					chromecast: c.chromecast && window.chrome && window.chrome.cast,
					p2pAppId: c.p2pAppId,
					highlight: c.highlight || [],
					video: {
						url: c.urls,
						type: type,
						pic: c.pics,
						customType: amplayer.player.customType
					},
					touchVideoChangeProgress: c.touchVideoChangeProgress,
					pluginOptions: {}
				};
				switch (type) {
					case 'flv':
						opts.pluginOptions.flv = {
							//mediaDataSource:{},
							config: {
								isLive: c.live,
								autoCleanupSourceBuffer: true,
								autoCleanupMinBackwardDuration: 60
							}
						}; break;
					case 'hls':
						opts.pluginOptions.hls = {
							enableWorker: true,
							liveBackBufferLength: 15,
							backBufferLength: 15,
							liveMaxBackBufferLength: 15,
							maxBufferSize: 0,
							maxBufferLength: 10,
							liveSyncDurationCount: 1,
						}; break;
				}
				opts.pluginOptions = $.extend(opts.pluginOptions, c.pluginOptions || {});
				var player = new DPlayer(opts);
				bindListener(player);
				$(amplayer.elemPrefix() + '#video').data('player', player);
				return player;
			},
			'dplayer': function (options) {
				var c = $.extend(amplayer.options, options || {});
				var ctn = amplayer.elemPrefix(true);
				var elem = ctn && $(ctn).length>0 ? $(ctn)[0] : null;
				if (!elem) elem = document.getElementById('video');
				var player = new DPlayer({
					container: elem,
					autoplay: c.autoPlay,
					live: c.live,
					logo: c.logo,
					screenshot: c.screenshot,
					airplay: c.airplay,
					chromecast: c.chromecast && window.chrome && window.chrome.cast,
					p2pAppId: c.p2pAppId,
					highlight: c.highlight || [],
					touchVideoChangeProgress: c.touchVideoChangeProgress,
					video: {
						url: c.urls,
						pic: c.pics
					}
				});
				bindListener(player);
				$(amplayer.elemPrefix() + '#video').data('player', player);
				return player;
			}
		},
		'cookie': {
			'data': {},
			'set': function (name, value, days) {
				var exp = new Date();
				exp.setTime(exp.getTime() + days * 24 * 60 * 60 * 1000);
				var cookie = name + '=' + escape(value) + ';path=' + win.location.pathname + ';expires=' + exp.toUTCString() + ';sameSite=Lax';
				if (amplayer.secure) cookie += ';secure=true';
				document.cookie = cookie;
				amplayer.cookie.data[name] = value;
			},
			'get': function (name) {
				if (typeof (amplayer.cookie.data[name]) != 'undefined') return amplayer.cookie.data[name];
				var arr = document.cookie.match(new RegExp('(^| )' + name + '=([^;]*)(;|$)'));
				if (arr != null) return unescape(arr[2]);
			},
			'remove': function (name) {
				amplayer.cookie.set(name, '', -10);
			}
		},
		'localStorage': {
			'data': {},
			'set': function (name, value, days) {
				if (value === undefined) { return amplayer.localStorage.remove(name) }
				win.localStorage.setItem(name, value);
				amplayer.localStorage.data[name] = value;
			},
			'get': function (name) {
				if (typeof (amplayer.localStorage.data[name]) != 'undefined') return amplayer.localStorage.data[name];
				return win.localStorage.getItem(name);
			},
			'remove': function (name) {
				win.localStorage.removeItem(name);
			}
		},
		'store': {
			'set': function (name, value, days) {
				return storer().set(name, value, days);
			},
			'get': function (name) {
				return storer().get(name);
			},
			'remove': function (name) {
				return storer().remove(name);
			}
		},
		'jump': function (jump) {
			if (!jump) return;
			if ($.isFunction(jump)) {
				return jump();
			}
			top.location.href = jump;
		}
	};

	var supportedLocalStorage = null;
	function isLocalStorageNameSupported() {
		if (supportedLocalStorage !== null) {
			return supportedLocalStorage;
		}
		try { supportedLocalStorage = ('localStorage' in win) }
		catch (err) { supportedLocalStorage = false }
		return supportedLocalStorage;
	}

	function storer() {
		if (isLocalStorageNameSupported()) return amplayer.localStorage;
		return amplayer.cookie;
	}

	function callListener(name, thisObj, args) {
		if (!amplayer.options.listeners) return;
		if (typeof amplayer.options.listeners[name] != 'function') return;
		amplayer.options.listeners[name].call(thisObj, args);
	}

	function bindListener(player) {
		player.on('loadstart', function () {
			var $video = $(player.video);
			$video.attr('playsinline', 'true');
			$video.attr('x5-playsinline', 'true');
			$video.attr('webkit-playsinline', 'true');
			if (player.video.paused && !$video.hasClass('dplayer-mobile')) $(amplayer.elemPrefix() + '.amplayer-center-button').show();
			callListener('loadstart', this, arguments)
		});
		player.on('loadeddata', function () {
			onLoaded(amplayer.options, player);
			callListener('loadeddata', this, arguments)
		});
		player.on('timeupdate', function () {
			applyFilmRange(amplayer.options, player, player.video.currentTime);
			callListener('timeupdate', this, arguments)
		});
		player.on('ended', function () {
			callListener('ended', this, arguments)
			amplayer.jump(amplayer.options.jump);
		});
		player.on('pause', function () {
			var $video = $(player.video);
			if (!$video.hasClass('dplayer-mobile')) $(amplayer.elemPrefix() + '.amplayer-center-button').show();
			callListener('pause', this, arguments)
		});
		player.on('play', function () {
			var $video = $(player.video);
			if (!$video.hasClass('dplayer-mobile')) $(amplayer.elemPrefix() + '.amplayer-center-button').hide();
			callListener('play', this, arguments)
		});
		player.on('error', function () {
			//console.dir(arguments)
			callListener('error', this, arguments)
		});
		$(amplayer.elemPrefix() + '.amplayer-play').click(function () {
			player.play();
		});
		$(amplayer.elemPrefix() + '.amplayer-backward').click(function () {
			let t = Math.max(player.video.currentTime - 10, 0);
			player.seek(t);
			player.controller.setAutoHide();
		});
		$(amplayer.elemPrefix() + '.amplayer-forward').click(function () {
			let t = Math.min(player.video.currentTime + 10, player.video.duration);
			player.seek(t);
			player.controller.setAutoHide();
		});
	}

	function onLoaded(play, player) {
		if (play.live) return;
		if (play.autoSkip && play.seek <= 0 && play.filmRange && play.filmRange.playRange.min > 0) {
			play.seek = play.filmRange.playRange.min;
		}
		var lastTime = play.take ? amplayer.store.get(play.take) : 0;
		var current = player.video.currentTime;
		//console.log(play.seek,lastTime);
		//console.dir(amplayer.options)
		if (!lastTime) {
			if (current == play.seek) return;
			return player.seek(play.seek);
		}
		if (player.video.duration - lastTime < 60 || play.seek > lastTime) {
			if (current == play.seek) return;
			player.seek(play.seek);
		} else {
			if (current == lastTime) return;
			player.seek(lastTime);
		}
	}

	function applyFilmRange(play, player, current) {
		//console.log(play.trys,current);
		if (play.trys > 0 && current > play.trys) {
			$(player.video).trigger('endoftrial', play, player, current); //试看结束
			//player.seek(0);
			player.notice('试看结束');
			return player.pause(); // 试看
		}
		if (play.live) return; // 直播模式
		if (play.take) amplayer.store.set(play.take, current, 30);
		if (play.filmRange && (play.autoSkip || play.autoSkipAd)) {
			var rg = play.filmRange;
			if (play.autoSkip) { // 跳过片头和片尾
				if (rg.playRange.min > 0 && current < rg.playRange.min) {
					return player.seek(rg.playRange.min);
				}
				var playRangeMax = rg.playRange.max;
				if (playRangeMax < 0 && (playRangeMax * -1) < player.video.duration) playRangeMax = player.video.duration + playRangeMax;
				if (playRangeMax > 0 && current >= playRangeMax) {
					return amplayer.jump(play.jump);
				}
			}
			if (play.autoSkipAd && rg.skipRange) { // 跳过广告
				for (var i = 0; rg.skipRange.length; i++) {
					var v = rg.skipRange[i];
					if (current >= v.min) {
						if (v.max < 0 && (v.max * -1) < player.video.duration) v.max = player.video.duration + v.max;
						if (current < v.max) {
							return player.seek(v.max);
						}
					}
				}
			}
		}
	}

	win.amplayer = amplayer;

})(window);