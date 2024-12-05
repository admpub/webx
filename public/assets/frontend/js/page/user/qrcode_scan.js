$(function(){
var resultContainer = document.getElementById('qr-reader-results'),lastDecodedText='',t=null;
function clearTimeoutT(){
  if(!t) return; 
  clearTimeout(t);t=null;
}
function onScanSuccess(decodedText) {
  if(!decodedText||lastDecodedText===decodedText) return;
  lastDecodedText=decodedText;
  clearTimeoutT();
  t=setTimeout(function(){lastDecodedText=''},5000);
  if(!String(decodedText).startsWith('coscms:')){
    resultContainer.innerText=decodedText
    $('#btn-copy-qrcode.hidden').removeClass('hidden');
    App.message({text:App.t('识别成功'),type:'success'});
    return;
  }
  $.post(FRONTEND_URL+'/user/qrcode/scan',{data:decodedText},function(r){
      if(r.Code!=1){
        resultContainer.innerText=r.Info
        $('#btn-copy-qrcode:not(.hidden)').addClass('hidden')
        lastDecodedText='';
        clearTimeoutT();
        return App.message({text:r.Info,type:'error'});
      }
      resultContainer.innerText=r.Info
      $('#btn-copy-qrcode:not(.hidden)').addClass('hidden')
      return App.message({text:r.Info,type:'success'});
  },'json').error(function(){
    lastDecodedText='';
    clearTimeoutT();
  });
}
function onScanFailure(error) { 
  //lastDecodedText='';
}
// Square QR box with edge size = 70% of the smaller edge of the viewfinder.
function qrboxFunction(viewfinderWidth, viewfinderHeight) {
  var minEdgeSize = Math.min(viewfinderWidth, viewfinderHeight);
  var qrboxSize = Math.floor(minEdgeSize * 0.7);
  qrboxSize = Math.min(250,qrboxSize);
  return {width: qrboxSize,height: qrboxSize};
}
var scanner = new Html5QrcodeScanner("qr-reader", {fps: 10, qrbox: qrboxFunction, rememberLastUsedCamera: true});
scanner.render(onScanSuccess,onScanFailure);
attachCopy('#btn-copy-qrcode');
})
