$(function(){
function setPageContentEditor(v){
    var editor=$('#pageContent').data('codemirror');
    if(!editor){
        window.setTimeout(function(){setPageContentEditor(v)},500);
        return;
    }
    var mime=type2mime(v);
    editor.setOption("mode", mime);
    editor.refresh();
}
function type2mime(v){
    var mime='';
    switch(v){
        case 'html':mime='text/html';break;
        case 'json':mime='text/javascript';break;
        case 'xml':mime='text/xml';break;
        //case 'redirect':mime='message/http';break;
        default:mime='null';
    }
    return mime;
}
var checkedPageType=$('#routePageEditForm input[name=pageType]:checked');
App.editor.codemirror('#pageVars', {
    lineNumbers: true,
    theme: 'ambiance',
    mode:  "application/javascript"
},'json');
App.editor.codemirror('#pageContent', {
    lineNumbers: true,
    theme: 'ambiance',
    mode:  type2mime(checkedPageType.val()),
});
$('#routePageEditForm input[name=pageType]').on('click',function(){
    var v=$(this).val(),$form=$('#routePageEditForm');
    $form.find('.pageType-hide-'+v).hide();
    $form.find('.pageType-show-'+v).show();
    $form.find('[class^="pageType-hide"]:not(.pageType-hide-'+v+')').show();
    setPageContentEditor(v);
});
checkedPageType.trigger('click');
});