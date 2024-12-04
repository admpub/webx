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
var html5QrcodeScanner = new Html5QrcodeScanner("qr-reader", {fps: 10, qrbox: 250, rememberLastUsedCamera: true});
html5QrcodeScanner.render(onScanSuccess,onScanFailure);
attachCopy('#btn-copy-qrcode');
})
