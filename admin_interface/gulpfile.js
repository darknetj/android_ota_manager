var gulp = require('gulp'),
    browserify = require('gulp-browserify'),
    browserSync = require('browser-sync'),
    concat = require('gulp-concat'),
    hbsfy = require('hbsfy'),
    connect = require('gulp-connect');

gulp.task('browserify', function() {
  gulp.src('src/js/main.js')
    .pipe(browserify({transform: 'hbsfy'}))
    .pipe(concat('main.js'))
    .pipe(gulp.dest('static/js'));
});

gulp.task('copy', function() {
  gulp.src('src/index.html')
    .pipe(gulp.dest('static'));
});

gulp.task('watch', ['build'], function() {
  gulp.watch('src/**/*.*', ['build', 'connectStop', 'connect']);
});

gulp.task('connect', function() {
  connect.server({
    root: 'static',
    port: 9000
  });
});

gulp.task('connectStop', function() {
  connect.serverClose();
});

gulp.task('default', ['watch','connect']);
gulp.task('build', ['browserify', 'copy']);
