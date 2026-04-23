$(function(){
App.showRequriedInputStar();
$('input[name="thirdparty[ip2region][value][mode]"]').on('click',function(){
    var v=$(this).val();
    if(v=='api'){
        $('#thirdparty-ip2region-api-settings').removeClass('hide');
        $('#thirdparty-ip2region-api-settings').find('[required]').prop('readonly',false);
    }else{
        $('#thirdparty-ip2region-api-settings:not(.hide)').addClass('hide');
        $('#thirdparty-ip2region-api-settings').find('[required]').prop('readonly',true);
    }
});
$('input[name="thirdparty[segment][value][engine]"]').on('click',function(){
    var v=$(this).val();
    if(v=='api'){
        $('#thirdparty-segment-settings').removeClass('hide');
        $('#thirdparty-segment-settings').find('[required]').prop('readonly',false);
    }else{
        $('#thirdparty-segment-settings:not(.hide)').addClass('hide');
        $('#thirdparty-segment-settings').find('[required]').prop('readonly',true);
    }
});
$('input[name="thirdparty[translate][value][on]"]').on('click',function(){
    var v=$(this).val();
    if(v=='1'){
        $('#thirdparty-translate-settings').removeClass('hide');
        $('select[name="thirdparty[translate][value][provider]"]').trigger('change');
    }else{
        $('#thirdparty-translate-settings:not(.hide)').addClass('hide');
        $('#thirdparty-translate-provider-config').find('[required]').prop('readonly',true);
    }
});
$('select[name="thirdparty[translate][value][provider]"]').on('change',function(){
    var v=$(this).val(),formElements=$(this).children(':selected').data('form');
    var b=$('#thirdparty-translate-provider-config'),formData=b.data('config');
    var h=template('tpl-thirdparty-translate-provider-forms', {formElements:formElements,formData:formData,required:'required'});
    b.html(h);
}).trigger('change');
$('input[name="thirdparty[segment][value][engine]"]:checked').trigger('click');
$('input[name="thirdparty[ip2region][value][mode]"]:checked').trigger('click');
$('input[name="thirdparty[translate][value][on]"]:checked').trigger('click');
})