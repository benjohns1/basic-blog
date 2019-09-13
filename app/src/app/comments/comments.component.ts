import { Component, OnInit, Input } from '@angular/core';
import { Post } from '../post';
import { NewComment, Comment } from '../comment';
import { BlogService } from '../blog.service';
import { NgForm } from '@angular/forms';

@Component({
  selector: 'app-comments',
  templateUrl: './comments.component.html',
  styleUrls: ['./comments.component.less']
})
export class CommentsComponent implements OnInit {

  @Input() public post: Post;

  public comment: NewComment = {
    postId: 0,
    commenter: "",
    body: "",
  };

  constructor(private blogService: BlogService) { }

  ngOnInit() {
  }

  public submitComment(form: NgForm) {
    this.comment.postId = this.post.id;
    this.blogService.newComment(this.comment).subscribe(() => {
      const comment: Comment = {
        id: 0,
        commenter: this.comment.commenter,
        body: this.comment.body,
        createdTime: new Date(),
      };
      this.post.comments.unshift(comment);
      form.reset();
    }, error => {
      // @TODO: show friendly error in UI
      console.error(error)
    });
  }

}
