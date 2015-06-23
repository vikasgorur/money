var App = React.createClass({displayName: "App",
  getInitialState: function() {
    return {outputText: ""};
  },

  onInputChange: function(event, value) {
    console.log("input was changed to " + value);
    var self = this;

    $.ajax("/convert", {
      data: {text: value},
      dataType: "json",
      success: function(data, status) {
        self.setState({outputText: data.text});
      },
      error: function(xhr, status) {
        self.setState({outputText: status});
      }
    });
  },

	render: function() {
    return (
      React.createElement("div", {className: "app"}, 
        React.createElement(Input, {onChange: this.onInputChange}), 
        React.createElement(Output, {text: this.state.outputText})
      )
    );
	}
});

// Contains the input box. Takes a callback prop that is
// called when the input changes, along with the new value.
var Input = React.createClass({displayName: "Input",
  onChange: function(event) {
    event.preventDefault();
    this.props.onChange(event, event.target.elements[0].value);
  },

  render: function() {
    return (
      React.createElement("div", {className: "input"}, 
        React.createElement("form", {onSubmit: this.onChange}, 
          React.createElement("input", {className: "input-box", type: "text", autoFocus: "true"})
        )
      )
    );
  }
});

var Output = React.createClass({displayName: "Output",
  render: function() {
    return (
      React.createElement("div", {className: "output"}, 
        React.createElement("span", {className: "output-box"}, this.props.text)
      )
    );
  }
});