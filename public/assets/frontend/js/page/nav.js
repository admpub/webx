$(function () {
	'use strict'
	var isNavSpy = $('body').hasClass('landing-3')||$('body').hasClass('corporate-3'),
		isNavSuperfish = $('#navbarMain').hasClass('nav-superfish'),
		hideClass = $('#navbarMain').hasClass('hidden-md-down') ? 'hidden-md-down' : 'hidden-sm-down';

	if (isNavSuperfish) {
		$('#navbarMain').superfish();
	}
	// animated smooth scroll on target from top menu
	$('#navbarMain a[href^="#"]').on('click', function(e){
	  //auto hide menu, mobile only
	  $('#navbarMain').addClass(hideClass);
	  var target = $(this).attr('href');
	  $('html, body').animate({scrollTop: $(''+target).offset().top}, 500);
	});
	// show hide menu, for mobile only
	$('#showMenu').on('click', function () {
		if ($('#navbarMain').hasClass(hideClass)) {
			$('#navbarMain').removeClass(hideClass);
			if(isNavSuperfish) $('body').addClass('show-sf-menu');
			if(isNavSpy) $('#navbarMain').addClass('nav-mobile');
			// change navicon icon to close icon button
			$(this).find('.fa').removeClass('fa-navicon').addClass('fa-close');
		} else {
			$('#navbarMain').addClass(hideClass);
			if(isNavSuperfish) $('body').removeClass('show-sf-menu');
			if(isNavSpy) $('#navbarMain').removeClass('nav-mobile');
			// change back from close to navicon icon button
			$(this).find('.fa').removeClass('fa-close').addClass('fa-navicon');
		}
		return false;
	});
});