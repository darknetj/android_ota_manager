var $                   = require('jquery'),
    Backbone            = require('backbone'),
    template            = require("../templates/layout_nav.hbs");

Backbone.$ = $;

module.exports = Backbone.View.extend({
    el: $('#nav'),
    template: template,

    initialize: function () {
        this.render();
    },

    render: function () {
        this.$el.html(template());
        return this;
    }
});
