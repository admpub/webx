$(function(){
function setPageContentEditor(elem,v){
    var editor=$(elem).data('codemirror');
    if(!editor){
        window.setTimeout(function(){setPageContentEditor(v)},500);
        return;
    }
    var mime=type2mime(v);
    editor.setOption("mode", mime);
    editor.refresh();
    //console.log('refresh',mime);
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
function initCodeMirror(elem){
    App.editor.codemirror(elem, { lineNumbers: true, theme: 'ambiance', mode:  type2mime(checkedPageType.val()) });
}
function setAllEditor(cb){
    $('#routePageEditForm').find('textarea[name$="[pageContent]"]').each(function(){
        if(this.id=='pageContent') return;
        cb('#'+this.id);
    })
    cb('#pageContent');
}
setAllEditor(initCodeMirror);
$('#routePageEditForm input[name=pageType]').on('click',function(){
    var v=$(this).val(),$form=$('#routePageEditForm');
    $form.find('.pageType-hide-'+v).hide();
    $form.find('.pageType-show-'+v).show();
    $form.find('[class^="pageType-hide"]:not(.pageType-hide-'+v+')').show();
    setAllEditor(function(elem){
        setPageContentEditor(elem,v);
    })
});
$('#routePageEditForm').find('.langset > .nav-tabs > li').on('click',function(){
    var tabContentE=$(this).children('a').attr('href');
    var editor=$(tabContentE).find('textarea[name$="[pageContent]"]').data('codemirror');
    if(!editor) return;
    setTimeout(function(){editor.refresh();},200);
})
checkedPageType.trigger('click');
});