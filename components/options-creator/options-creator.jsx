import React from "react";
import ReactDOM from "react-dom";

class Option extends React.Component {
  constructor() {
    super();
    this.setName = this.setName.bind(this);
    this.setLang = this.setLang.bind(this);
    this.setChoices = this.setChoices.bind(this);
    this.setStocks = this.setStocks.bind(this);
    this.remove = this.remove.bind(this);
    this.state = {
      value: ":: ::es::"
    };
  }

  setStocks(stocks) {
    this.setState({
      value: `${this.state.value.split("::")[0]}::${
        this.state.value.split("::")[1]
      }::${this.state.value.split("::")[2]}::${stocks}`
    });
  }

  setName(name) {
    this.setState({
      value: `${name}::${this.state.value.split("::")[1]}::${
        this.state.value.split("::")[2]
      }::${this.state.value.split("::")[3]}`
    });
  }

  setLang(lang) {
    this.setState({
      value: `${this.state.value.split("::")[0]}::${
        this.state.value.split("::")[1]
      }::${lang}::${this.state.value.split("::")[3]}`
    });
  }

  setChoices(desc) {
    this.setState({
      value: `${this.state.value.split("::")[0]}::${desc}::${
        this.state.value.split("::")[2]
      }::${this.state.value.split("::")[3]}`
    });
  }

  remove() {
    this.me.parentNode.removeChild(this.me);
  }

  render() {
    return (
      <div ref={el => (this.me = el)}>
        <input
          onChange={e => this.setName(e.target.value)}
          className="uk-input uk-width-1-6"
          style={{ margin: "10px" }}
          type="text"
          placeholder="Name"
        />
        <select
          onChange={e => this.setLang(e.target.value)}
          class="uk-select uk-width-1-6"
        >
          <option value="es">es</option>
          <option value="en">en</option>
        </select>
        <input
          onChange={e => this.setChoices(e.target.value)}
          className="uk-input uk-width-2-3"
          style={{ margin: "10px" }}
          type="text"
          placeholder="Description"
        />
        <input
          onChange={e => this.setStocks(e.target.value)}
          className="uk-input uk-width-2-3"
          style={{ margin: "10px" }}
          type="text"
          placeholder="Stock"
        />
        <a onClick={this.remove} uk-icon="icon: close" />
        <input type="hidden" name="Options[]" value={this.state.value} />
      </div>
    );
  }
}

const NewOption = ({ handler }) => (
  <button className="uk-button uk-button-default" onClick={handler}>
    New Option
  </button>
);

class Creator extends React.Component {
  constructor() {
    super();
    this.new = this.new.bind(this);
    this.elements = [<NewOption handler={this.new} />];
  }

  new(event) {
    event.preventDefault();
    event.stopPropagation();
    this.elements.push(<Option />);
    this.forceUpdate();
  }

  render() {
    return <div ref={el => (this.el = el)}>{this.elements.map(el => el)}</div>;
  }
}

ReactDOM.render(<Creator />, document.getElementById("options-creator"));
