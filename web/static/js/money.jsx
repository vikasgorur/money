var App = React.createClass({
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
      <div className="app">
        <Input onChange={this.onInputChange} />
        <Output text={this.state.outputText} />
      </div>
    );
	}
});

// Contains the input box. Takes a callback prop that is
// called when the input changes, along with the new value.
var Input = React.createClass({
  onChange: function(event) {
    event.preventDefault();
    this.props.onChange(event, event.target.elements[0].value);
  },

  render: function() {
    return (
      <div className="input">
        <form onSubmit={this.onChange}>
          <input className="input-box" type="text" autoFocus="true" />
        </form>
      </div>
    );
  }
});

var Output = React.createClass({
  render: function() {
    return (
      <div className="output">
        <span className="output-box">{this.props.text}</span>
      </div>
    );
  }
});