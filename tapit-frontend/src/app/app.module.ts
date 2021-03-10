import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MainComponent } from './main/main.component';
import { CampaignComponent } from './campaign/campaign.component';
import { CampaignNewComponent } from './campaign-new/campaign-new.component';
import { NotificationComponent } from './notification/notification.component';
import { PhonebookComponent } from './phonebook/phonebook.component';
import { PhonebookNewComponent } from './phonebook-new/phonebook-new.component';
import { TextTemplateComponent } from './text-template/text-template.component';
import { TextTemplateNewComponent } from './text-template-new/text-template-new.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { ProviderComponent } from './provider/provider.component';
import { ProfileComponent } from './profile/profile.component';
import { CampaignViewComponent } from './campaign-view/campaign-view.component';
import { WebTemplateComponent } from './web-template/web-template.component';
import { WebTemplateNewComponent } from './web-template-new/web-template-new.component';
import { GlobalSettingsComponent } from './global-settings/global-settings.component';

@NgModule({
  declarations: [
    AppComponent,
    MainComponent,
    CampaignComponent,
    CampaignNewComponent,
    NotificationComponent,
    PhonebookComponent,
    PhonebookNewComponent,
    TextTemplateComponent,
    TextTemplateNewComponent,
    LoginComponent,
    RegisterComponent,
    ProviderComponent,
    ProfileComponent,
    CampaignViewComponent,
    WebTemplateComponent,
    WebTemplateNewComponent,
    GlobalSettingsComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    HttpClientModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
