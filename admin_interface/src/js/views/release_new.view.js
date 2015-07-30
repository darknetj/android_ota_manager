var $                   = require('jquery'),
    Backbone            = require('backbone'),
    Release             = require("../models/release.model");
    template            = require("../templates/release_new.hbs");

Backbone.$ = $;

module.exports = Backbone.View.extend({

    initialize: function () {
        this.render();
    },

    render: function () {
        this.$el.html(template(this.model.attributes));
        this.input = this.$('.edit');
        return this;
    }

    events: {
        "keypress form#release":  "createOnEnter",
        "click #new-release": "create",
    },

    create: function(e) {
        if (!this.input.val()) return;
        Release.create({title: this.input.val()});
    },

    createOnEnter: function(e) {
        if (e.keyCode != 13) return;
        create();
    }
});
