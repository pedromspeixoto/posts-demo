import React from "react";

class Posts extends React.Component {
    state = {
        posts: []
    }
    componentDidMount() {
        let API_URL = process.env.REACT_APP_API_URL ? process.env.REACT_APP_API_URL : "http://localhost:8080"
        // Declare API_URL as the value of the REACT_APP_API_URL environment variable, if not defined, static value
        if (window._env_ !== undefined && window._env_.REACT_APP_API_URL !== undefined) {
            API_URL = window._env_.REACT_APP_API_URL
        }
        fetch(API_URL + "/v1/posts")
        .then(res => res.json())
        .then((data) => {
            this.setState({ posts: data.data.data.posts })
        })
        .catch(console.log)
    }
    render() {
        return (
            <div>
                <div className="row mb-3">
                    <div className="col-md-12 text-md-center">
                        <h2 className="text-poppins font-weight-bold">Posts</h2>
                    </div>
                </div>
                <div className="row justify-content-center">
                    {this.state.posts.map((post) => (
                        <div className="col-md-8">
                        <div className="card card-outline translate-up mb-3">
                            <div className="card-body">
                                <h5 className="card-title">#<strong className="text-primary">{post.post_id}</strong></h5>
                                <h6 className="body-1 text-muted">{post.content}</h6>
                            </div>
                            <div className="card-footer text-right">
                                <h6 className="subtitle-1 text-muted">{post.created_at}</h6>
                            </div>
                        </div>
                    </div>
                    ))}
                </div>
            </div>
        );
    }
}

export default Posts;