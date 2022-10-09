$(function(){
  // animated smooth scroll on target from top menu
  $('#navbarMain .nav-link').on('click', function(e){
    //auto hide menu, mobile only
    $('#navbarMain').addClass('hidden-sm-down');
    var target = $(this).attr('href');
    $('html, body').animate({scrollTop: $(''+target).offset().top}, 500);
  });
  $('#showMenu').on('click', function(){
    if($('#navbarMain').hasClass('hidden-sm-down')) {
      $('#navbarMain').removeClass('hidden-sm-down');
      $('#navbarMain').addClass('nav-mobile');
      // change navicon icon to close icon button
      $(this).find('.fa').removeClass('fa-navicon').addClass('fa-close');
    } else {
      $('#navbarMain').addClass('hidden-sm-down');
      $('#navbarMain').removeClass('nav-mobile');
      // change back from close to navicon icon button
      $(this).find('.fa').removeClass('fa-close').addClass('fa-navicon');
    }
    return false;
  });
});