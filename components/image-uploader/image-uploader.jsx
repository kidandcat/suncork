import React from "react";
import ReactDOM from "react-dom";

class Image extends React.Component {
  constructor(props) {
    super(props);
    this.handleChange = this.handleChange.bind(this);
    this.state = {
      inputMode: true,
      image: ""
    };
  }
  handleChange(files) {
    if (files && files[0]) {
      const reader = new FileReader();

      reader.onload = e => {
        if(this.state.image == ""){
          this.props.trigger();
        }
        this.setState({
          image: e.target.result,
          inputMode: false
        });
      };

      reader.readAsDataURL(files[0]);
    }
  }

  render() {
    return (
      <div style={{
        display: "inline-block"
      }}>
        <input
          /*style={{
            opacity: this.state.inputMode ? 1 : 0,
            position: this.state.inputMode ? "static" : "absolute"
          }}*/
          onChange={e => this.handleChange(e.target.files)}
          type="file"
          name="Image"
          class="uk-input uk-width-1-2"
          style={{
            display: "block"
          }}
        />

        <input value={this.state.image} name="Image" type="hidden" />

        <img
          src={this.state.image}
          style={{ opacity: this.state.inputMode ? 0 : 1, maxWidth: "250px" }}
        />
      </div>
    );
  }
}

class Uploader extends React.Component {
  constructor() {
    super();
    this.new = this.new.bind(this);
    this.elements = [<Image trigger={this.new} />];
  }

  new() {
    this.elements.push(<Image trigger={this.new} />);
    this.forceUpdate();
  }

  render() {
    return <div ref={el => (this.el = el)}>{this.elements.map(el => el)}</div>;
  }
}

ReactDOM.render(<Uploader />, document.getElementById("image-uploader"));
